package couchcampaign

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
)

// Game options.
const (
	// The number of cards drawn from a single player's deck between elections.
	//
	// Once at least one player draws this many cards during an offseason, the
	// next campaign season begins. Once all players have drawn at least this
	// many cards, voting season begins. Voting season ends when all clients are
	// finished counting votes and the election results are released.
	//
	// This offseason -> campaign -> voting scheme is used to reward players who
	// swipe cards quickly and avoid waiting on slower players.
	numCardsBetweenElections = 10
)

var welcomeCard = infoCard{`
	Welcome to couchcampaign!

	Today is your first day in office as the governor.
	Work with your advisors to stay in office as long as you can.
	Remember, the ultimate goal is to win the presidency.
`}

type Game struct {
	election *electionStateMachine
	pids     []PID
	clients  map[PID]*ClientWorker
	jobs     map[PID]chan ClientJob
	state    map[PID]playerState
	decks    map[PID]*Deck
	inputs   chan ClientMessage
	errors   chan ClientError

	baseCards []Card

	mu sync.Mutex
}

func NewGame(clients map[PID]*ClientWorker) (*Game, error) {
	g := &Game{
		clients: clients,
		pids:    make([]PID, 0, len(clients)),
		decks:   make(map[PID]*Deck),
		state:   make(map[PID]playerState),
		jobs:    make(map[PID]chan ClientJob),
		inputs:  make(chan ClientMessage, len(clients)),
		errors:  make(chan ClientError, len(clients)),
	}

	for pid := range clients {
		g.pids = append(g.pids, pid)
	}
	g.election = newElectionStateMachine(numCardsBetweenElections, g.pids)

	var err error
	g.baseCards, err = loadBaseCards()
	if err != nil {
		return nil, fmt.Errorf("loadBaseCards: %w", err)
	}

	for _, pid := range g.pids {
		g.decks[pid] = g.buildBaseDeck()
		g.jobs[pid] = make(chan ClientJob, 2)
		g.state[pid] = newPlayerState()
		// Insert the starting card for all players.
		g.decks[pid].InsertCardWithPriority(welcomeCard, maxCardPriority)
	}

	return g, nil
}

func (g *Game) Run(_ context.Context) {
	defer g.shutdown()

	// Spawn all clients.
	for _, pid := range g.pids {
		go g.clients[pid].Run(g.jobs[pid], g.inputs, g.errors)
		g.sendNextJob(pid)
	}

	for !g.isOver() {
		select {
		case input := <-g.inputs:
			g.handleInput(input)
		case err := <-g.errors:
			g.handleError(err)
		}
	}

	log.Println("Game over")
}

func (g *Game) shutdown() {
	close(g.inputs)
	close(g.errors)
	for _, pid := range g.pids {
		g.disconnect(pid)
	}
}

func (g *Game) disconnect(pid PID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.clients, pid)
	close(g.jobs[pid])
	delete(g.jobs, pid)
	delete(g.state, pid)
	for i, id := range g.pids {
		if id == pid {
			g.pids = append(g.pids[i:], g.pids[:i+1]...)
		}
	}
}

func (g *Game) handleError(err ClientError) {
	if IsConnectionCloseError(err) {
		g.disconnect(err.PID)
		log.Printf("disconnected from client %v", err.PID)
		return
	}
	log.Printf("%v: %v", err.PID, err.Error())
}

func (g *Game) handleInput(input ClientMessage) {
	defer g.debugDumpDecks()

	state := g.state[input.PID]
	switch input.Input {
	case DismissInfoCardInput:
		g.decks[input.PID].RemoveCard(input.Card)
		g.state[input.PID] = OnDismissInfoCard(state, input.Card.(infoCard))
		g.sendNextJob(input.PID)
	case AcceptActionCardInput:
		g.decks[input.PID].RemoveCard(input.Card)
		g.state[input.PID] = OnAcceptActionCard(state, input.Card.(actionCard))
		if g.state[input.PID] == EmptyPlayerState {
			g.decks[input.PID].Clear()
			crumbleCard := infoCard{"Society has crumbled and you are being forced out of office."}
			g.decks[input.PID].InsertCardWithPriority(crumbleCard, maxCardPriority)
			g.decks[input.PID] = g.buildBaseDeck()
			g.state[input.PID] = newPlayerState()
		}
		g.updateElectionState(input)
		g.sendNextJob(input.PID)
	case RejectActionCardInput:
		g.decks[input.PID].RemoveCard(input.Card)
		g.state[input.PID] = OnRejectActionCard(state, input.Card.(actionCard))
		if g.state[input.PID] == EmptyPlayerState {
			g.decks[input.PID].Clear()
			crumbleCard := infoCard{"Society has crumbled and you are being forced out of office."}
			g.decks[input.PID].InsertCardWithPriority(crumbleCard, maxCardPriority)
			g.decks[input.PID] = g.buildBaseDeck()
			g.state[input.PID] = newPlayerState()
		}
		g.updateElectionState(input)
		g.sendNextJob(input.PID)
	default:
		if _, ok := input.Card.(votingCard); ok {
			g.updateElectionState(input)
		}
	}
}

func (g *Game) updateElectionState(input ClientMessage) {
	oldSeason := g.election.CurrentSeason()
	newSeason := g.election.HandleCardPlayed(input.PID, input.Card)
	if oldSeason == newSeason {
		return
	}

	switch {
	case oldSeason == offSeason && newSeason == campaignSeason:
		g.annouce("Campaign season has begun!")
		break
	case oldSeason == campaignSeason && newSeason == votingSeason:
		g.annouce("Voting season has begun!")
		for _, pid := range g.pids {
			g.decks[pid].InsertCard(theVotingCard)
		}
		break
	case oldSeason == votingSeason && newSeason == offSeason:
		// If we went from voting season to the off season then all players were
		// looking at the voting card. Pop it from every deck, announce the
		// results, and the move clients to the next card.
		for _, pid := range g.pids {
			g.decks[pid].RemoveCard(theVotingCard)
		}
		winner := g.computeElectionWinner()
		for _, pid := range g.pids {
			if pid == winner {
				g.decks[pid].InsertCardWithPriority(infoCard{"You won the election! Now get back to work."}, maxCardPriority)
			} else {
				g.decks[pid].InsertCardWithPriority(infoCard{"You lost the election. Now get back to to work."}, maxCardPriority)
			}
		}
		g.annouce("The offseason has begun!")
		for _, pid := range g.pids {
			g.sendNextJob(pid)
		}
		break
	default:
		log.Fatalf("invalid transition from %s to %s", oldSeason, newSeason)
	}
}

func (g *Game) annouce(message string) {
	for _, pid := range g.pids {
		g.decks[pid].InsertCard(infoCard{message})
	}
}

func (g *Game) buildBaseDeck() *Deck {
	cards := make([]Card, len(g.baseCards))
	copy(cards, g.baseCards)
	deck := NewDeck(cards)
	deck.ShuffleActionCards()
	return deck
}

func (g *Game) computeElectionWinner() PID {
	highscore := math.MinInt64
	scores := make(map[int][]PID)
	for _, pid := range g.pids {
		score := g.state[pid].SocietyScore()
		scores[score] = append(scores[score], pid)
		if score > highscore {
			highscore = score
		}
	}

	if len(scores[highscore]) == 1 {
		return scores[highscore][0]
	}

	// If there's a tie, whoever tied and has the most balanced stats wins.
	minVariance := math.MaxFloat64
	variances := make(map[float64][]PID)
	for _, pid := range scores[highscore] {
		variance := g.state[pid].SocietyVariance()
		variances[variance] = append(variances[variance], pid)
		if variance < minVariance {
			minVariance = variance
		}
	}

	if len(variances[minVariance]) == 1 {
		return variances[minVariance][0]
	}

	// If there's a tie, the winner is randomly chosen, just like real life.
	ties := variances[minVariance]
	return ties[rand.Intn(len(ties))]
}

func (g *Game) debugDumpDecks() {
	for _, pid := range g.pids {
		log.Printf("==== DECK %v ====\n", pid)
		g.decks[pid].DebugDump(os.Stderr)
		log.Printf("==== END DECK %v ====\n\n", pid)
	}
}

func (g *Game) isOver() bool {
	return len(g.clients) <= 0
}

func (g *Game) sendNextJob(pid PID) {
	if g.decks[pid].IsEmpty() {
		g.decks[pid] = g.buildBaseDeck()
	}
	nextCard := g.decks[pid].TopCard()

	g.jobs[pid] <- ClientJob{
		Card:  nextCard,
		Stats: g.state[pid],
	}
	log.Printf("Sent top card to %v", pid)
}

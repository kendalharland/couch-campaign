package couchcampaign

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/gobuffalo/uuid"
	"github.com/gorilla/websocket"
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

var welcomeMessage = `
	Welcome to couchcampaign!

	Today is your first day in office as the governor.
	Work with your advisors to stay in office as long as you can.
	Remember, the ultimate goal is to win the presidency.
`

type message struct {
	pid   uuid.UUID
	card  Card
	input string
}

type Game struct {
	election *electionStateMachine
	pids     []uuid.UUID
	clients  map[uuid.UUID]*Client
	jobs     map[uuid.UUID]chan func(*Client)
	stats    map[uuid.UUID]*stats
	decks    map[uuid.UUID]*Deck
	inputs   chan message

	baseCards []Card
}

func NewGame(clients map[uuid.UUID]*Client) (*Game, error) {
	pids := make([]uuid.UUID, 0, len(clients))
	for pid := range clients {
		pids = append(pids, pid)
	}

	baseCards, err := loadBaseCards()
	if err != nil {
		return nil, fmt.Errorf("loadBaseCards: %w", err)
	}
	log.Println(baseCards)

	g := &Game{
		election:  newElectionStateMachine(numCardsBetweenElections, pids),
		pids:      pids,
		clients:   clients,
		decks:     make(map[uuid.UUID]*Deck),
		jobs:      make(map[uuid.UUID]chan func(*Client)),
		stats:     make(map[uuid.UUID]*stats),
		inputs:    make(chan message, len(clients)),
		baseCards: baseCards,
	}

	for _, pid := range pids {
		g.decks[pid] = g.buildBaseDeck()
		g.jobs[pid] = make(chan func(*Client), 2)
		g.stats[pid] = newStats()
	}

	return g, nil
}

func (g *Game) Run(_ context.Context) {
	for _, pid := range g.pids {
		card := infoCard{welcomeMessage}
		g.decks[pid].InsertCardWithPriority(card, maxCardPriority)
	}

	for _, client := range g.clients {
		go client.run(g.jobs[client.pid])
		g.sendTopCard(client.pid)
	}

	for !g.isOver() {
		g.handleInput(<-g.inputs)
	}

	close(g.inputs)
	for _, c := range g.jobs {
		close(c)
	}

	log.Println("Game over")
}

func (g *Game) disconnect(pid uuid.UUID) {
	log.Printf("Client close: %v", g.clients[pid].Close())
	delete(g.clients, pid)
	delete(g.jobs, pid)
	delete(g.stats, pid)
	for i, id := range g.pids {
		if id == pid {
			g.pids = append(g.pids[i:], g.pids[:i+1]...)
		}
	}
}

func (g *Game) handleInput(input message) {
	defer g.debugDumpDecks()

	switch card := input.card.(type) {
	case infoCard:
		g.decks[input.pid].RemoveCard(input.card)
		g.sendTopCard(input.pid)
	case actionCard:
		s := g.stats[input.pid]
		switch input.input {
		case "accept":
			card.accept(s)
		case "reject":
			card.reject(s)
		default:
			log.Printf("invalid action: %q", input.input)
			return
		}
		switch {
		case s.Wealth <= minStatValue || maxStatValue <= s.Wealth:
			g.decks[input.pid].Clear()
			g.decks[input.pid].InsertCardWithPriority(infoCard{"Your state is bankrupt and you are being forced out of office."}, maxCardPriority)
			g.decks[input.pid] = g.buildBaseDeck()
			g.stats[input.pid] = newStats()
		case s.Health <= minStatValue || maxStatValue <= s.Health:
			g.decks[input.pid].Clear()
			g.decks[input.pid].InsertCardWithPriority(infoCard{"Basically everyone in your state is dead. you've been removed from office."}, maxCardPriority)
			g.decks[input.pid] = g.buildBaseDeck()
			g.stats[input.pid] = newStats()
		case s.Stability <= minStatValue || maxStatValue <= s.Stability:
			g.decks[input.pid].Clear()
			g.decks[input.pid].InsertCardWithPriority(infoCard{"People hate living here so they've staged a coup and removed you from office."}, maxCardPriority)
			g.decks[input.pid] = g.buildBaseDeck()
			g.stats[input.pid] = newStats()
		}
		g.decks[input.pid].RemoveCard(input.card)
		g.updateElectionState(input)
		g.sendTopCard(input.pid)
	case votingCard:
		g.updateElectionState(input)
	}
}

func (g *Game) updateElectionState(input message) {
	oldSeason := g.election.CurrentSeason()
	newSeason := g.election.HandleCardPlayed(input.pid, input.card)
	if oldSeason == newSeason {
		return
	}

	switch {
	case oldSeason == offSeason && newSeason == campaignSeason:
		g.annouce("Campaign season has begun!")
		break
	case oldSeason == campaignSeason && newSeason == votingSeason:
		g.annouce("Voting season has begun!")
		g.waitForVotes()
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
			g.sendTopCard(pid)
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

func (g *Game) computeElectionWinner() uuid.UUID {
	// Whoever has the highest average score wins.
	highscore := math.MinInt64
	scores := make(map[int][]uuid.UUID)
	for _, pid := range g.pids {
		score := g.stats[pid].Sum()
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
	variances := make(map[float64][]uuid.UUID)
	for _, pid := range scores[highscore] {
		variance := g.stats[pid].Variance()
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
	for pid := range g.clients {
		if !g.decks[pid].IsEmpty() {
			return false
		}
	}
	return len(g.clients) > 0
}

func (g *Game) sendTopCard(pid uuid.UUID) {
	if g.decks[pid].IsEmpty() {
		g.decks[pid] = g.buildBaseDeck()
	}

	c := g.decks[pid].TopCard()
	s := g.stats[pid]

	if c == theVotingCard {
		g.jobs[pid] <- func(n *Client) {
			n.showCard(theVotingCard, s)
			g.inputs <- message{
				pid:  n.pid,
				card: c,
			}
		}
		return
	}

	g.jobs[pid] <- func(cli *Client) {
		c := g.decks[cli.pid].TopCard()
		cli.showCard(c, s)

		input, err := cli.getInput()
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr) {
			g.disconnect(pid)
			return
		}
		if err != nil {
			log.Println(err)
			return
		}

		g.inputs <- message{
			pid:   cli.pid,
			card:  c,
			input: input,
		}
	}
}

func (g *Game) waitForVotes() {
	for _, pid := range g.pids {
		g.decks[pid].InsertCard(theVotingCard)
	}
}

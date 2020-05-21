package couchcampaign

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
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

type GameState struct {
	pids         []PID
	playerStates map[PID]playerState
	decks        map[PID]*Deck
	election     *electionStateMachine

	// This is a hack. delete.
	baseCards []Card
}

func newGameState(pids []PID, baseCards []Card) GameState {
	state := GameState{
		pids:         pids,
		baseCards:    baseCards,
		decks:        make(map[PID]*Deck),
		playerStates: make(map[PID]playerState),
		election:     newElectionStateMachine(numCardsBetweenElections, pids),
	}
	for _, pid := range pids {
		state.playerStates[pid] = newPlayerState()
		state.decks[pid] = state.buildBaseDeck()
		state.decks[pid].InsertCardWithPriority(welcomeCard, maxCardPriority)
	}
	return state
}

func (g *GameState) announce(message string) {
	for _, pid := range g.pids {
		g.decks[pid].InsertCard(infoCard{message})
	}
}

func (g *GameState) getNextJob(pid PID) ClientJob {
	if g.decks[pid].IsEmpty() {
		g.decks[pid] = g.buildBaseDeck()
	}
	nextCard := g.decks[pid].TopCard()

	return ClientJob{
		PID:   pid,
		Card:  nextCard,
		Stats: g.playerStates[pid],
	}
}

func (g *GameState) buildBaseDeck() *Deck {
	cards := make([]Card, len(g.baseCards))
	copy(cards, g.baseCards)
	deck := NewDeck(cards)
	deck.ShuffleActionCards()
	return deck
}

func (g *GameState) computeElectionWinner() PID {
	highscore := math.MinInt64
	scores := make(map[int][]PID)
	for _, pid := range g.pids {
		score := g.playerStates[pid].SocietyScore()
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
		variance := g.playerStates[pid].SocietyVariance()
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

func (g *GameState) removePlayer(pid PID) {
	delete(g.playerStates, pid)
}

func (g *GameState) handleMessage(input ClientMessage) (jobs []ClientJob, err error) {
	defer g.debugDumpDecks()

	state := g.playerStates[input.PID]
	switch input.Input {
	case DismissInfoCardInput:
		g.OnPlayerDismissedInfoCard(input.PID)
		jobs = []ClientJob{g.getNextJob(input.PID)}
	case AcceptActionCardInput:
		g.decks[input.PID].RemoveCard(input.Card)
		g.playerStates[input.PID] = OnAcceptActionCard(state, input.Card.(actionCard))
		if g.playerStates[input.PID] == EmptyPlayerState {
			g.decks[input.PID].Clear()
			crumbleCard := infoCard{"Society has crumbled and you are being forced out of office."}
			g.decks[input.PID].InsertCardWithPriority(crumbleCard, maxCardPriority)
			g.decks[input.PID] = g.buildBaseDeck()
			g.playerStates[input.PID] = newPlayerState()
		}
		jobs, err = g.updateElectionState(input)
	case RejectActionCardInput:
		g.decks[input.PID].RemoveCard(input.Card)
		g.playerStates[input.PID] = OnRejectActionCard(state, input.Card.(actionCard))
		if g.playerStates[input.PID] == EmptyPlayerState {
			g.decks[input.PID].Clear()
			crumbleCard := infoCard{"Society has crumbled and you are being forced out of office."}
			g.decks[input.PID].InsertCardWithPriority(crumbleCard, maxCardPriority)
			g.decks[input.PID] = g.buildBaseDeck()
			g.playerStates[input.PID] = newPlayerState()
		}
		jobs, err = g.updateElectionState(input)
	default:
		jobs, err = g.updateElectionState(input)
	}

	return jobs, err
}

func (g *GameState) debugDumpDecks() {
	for _, pid := range g.pids {
		log.Printf("==== DECK %v ====\n", pid)
		g.decks[pid].DebugDump(os.Stderr)
		log.Printf("==== END DECK %v ====\n\n", pid)
	}
}

func (g *GameState) updateElectionState(input ClientMessage) (jobs []ClientJob, err error) {
	oldSeason := g.election.CurrentSeason()
	newSeason := g.election.HandleCardPlayed(input.PID, input.Card)
	if oldSeason == newSeason {
		return jobs, nil
	}

	switch {
	case oldSeason == offSeason && newSeason == campaignSeason:
		g.announce("Campaign season has begun!")
		break
	case oldSeason == campaignSeason && newSeason == votingSeason:
		g.announce("Voting season has begun!")
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
		g.announce("The offseason has begun!")
		for _, pid := range g.pids {
			jobs = append(jobs, g.getNextJob(pid))
		}
		break
	default:
		return nil, fmt.Errorf("invalid transition from %s to %s", oldSeason, newSeason)
	}

	return jobs, nil
}

// OnPlayerDismissedInfoCard is called when an InfoCard is dismissed.
//
// pid is the PID of the player that dismissed the card.
// The player state represented by pid is the only player state that is updated.
func (g *GameState) OnPlayerDismissedInfoCard(pid PID) {
	ps := g.playerStates[pid]
	g.decks[pid].RemoveCard(ps.Card)
	ps = checkPlayerState(ps)
	g.playerStates[pid] = ps
}

func OnVotingCard(s playerState, _ votingCard) (next playerState) {
	next = s
	return checkPlayerState(next)
}

func OnAcceptActionCard(s playerState, c actionCard) (next playerState) {
	next.Wealth = s.Wealth + c.AccWealthEffect
	next.Health = s.Health + c.AccHealthEffect
	next.Stability = s.Stability + c.AccStabilityEffect
	return checkPlayerState(next)
}

func OnRejectActionCard(s playerState, c actionCard) (next playerState) {
	next.Wealth = s.Wealth + c.RejWealthEffect
	next.Health = s.Health + c.RejHealthEffect
	next.Stability = s.Stability + c.RejStabilityEffect
	return checkPlayerState(next)
}

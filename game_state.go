package cardgame

import (
	"log"
	"os"

	"google.golang.org/protobuf/proto"

	pb "couchcampaign/proto"
)

type GameState struct {
	pids      []PID
	cards     CardService
	decks     map[PID]*Deck
	societies map[PID]SocietyState
}

func NewGameState(pids []PID, cards CardService) *GameState {
	state := &GameState{
		pids:      pids,
		cards:     cards,
		decks:     make(map[PID]*Deck),
		societies: make(map[PID]SocietyState),
	}
	for _, pid := range pids {
		state.societies[pid] = newSocietyState()
		state.decks[pid] = cards.BuildIntroDeck()
	}
	return state
}

// OnCardAccepted is called when a Card is accepted.
//
// pid is the PID of the player that dismissed the card.
// Returns a list of PIDs for the players whose state was updated.
func (g *GameState) OnCardAccepted(pid PID) ([]PID, error) {
	society := g.societies[pid]
	g.decks[pid].RemoveCard(society.CardRef)
	g.cards.Card(society.CardRef).OnAccept(g, pid, society.CardRef)

	// Move this to game logic.
	society = checkSocietyState(society)
	if society == NilSocietyState {
		g.decks[pid] = g.cards.BuildSocietyCrumbledDeck()
		society = newSocietyState()
	}

	g.societies[pid] = society
	return nil, nil
}

// OnCardRejected is called when a Card is rejected.
//
// pid is the PID of the player that dismissed the card.
// Returns a list of PIDs for the players whose state was updated.
func (g *GameState) OnCardRejected(pid PID) ([]PID, error) {
	society := g.societies[pid]
	g.decks[pid].RemoveCard(society.CardRef)
	g.cards.Card(society.CardRef).OnReject(g, pid, society.CardRef)

	// Move this to game logic.
	society = checkSocietyState(society)
	if society == NilSocietyState {
		g.decks[pid] = g.cards.BuildSocietyCrumbledDeck()
		g.societies[pid] = newSocietyState()
	}

	g.societies[pid] = society
	return nil, nil
}

// OnCardShown is called when a card is shown to the player.
//
// pid is the PID of the player that saw the card.
// Returns a list of PIDs for the players whose state was updated.
func (g *GameState) OnCardShown(pid PID) ([]PID, error) {
	society := g.societies[pid]
	g.cards.Card(society.CardRef).OnShown(g, pid, society.CardRef)

	// Move this to game logic.
	society = checkSocietyState(society)
	return nil, nil
}

func (g *GameState) getNextJob(pid PID) ClientJob {
	if g.decks[pid].IsEmpty() {
		g.decks[pid] = g.cards.BuildBaseDeck()
	}

	card := g.cards.Card(g.decks[pid].TopCard())
	society := g.societies[pid]
	state, err := proto.Marshal(&pb.Message{
		Content: &pb.Message_PlayerState{
			PlayerState: &pb.PlayerState{
				Card:      card.toProto(),
				Leader:    society.Leader,
				Wealth:    int32(society.Wealth),
				Health:    int32(society.Health),
				Stability: int32(society.Stability),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return ClientJob{
		PID:           pid,
		State:         state,
		RequiresInput: g.cards.CardRequiresInput(card.ID),
	}
}

func (g *GameState) removePlayer(pid PID) {
	delete(g.societies, pid)
	delete(g.decks, pid)
}

func (g *GameState) debugDumpDecks() {
	for _, pid := range g.pids {
		log.Printf("==== DECK %v ====\n", pid)
		g.decks[pid].DebugDump(os.Stderr)
		log.Printf("==== END DECK %v ====\n\n", pid)
	}
}

func (g *GameState) DeckInsertAll(c CardRef) {
	for _, pid := range g.pids {
		g.decks[pid].InsertCard(c, g.cards.CardPriorityByType(c))
	}
}
func (g *GameState) DeckInsert(pid PID, card CardRef)          {}
func (g *GameState) DeckRemove(pid PID, card CardRef)          {}
func (g *GameState) DeckShuffle(pid PID)                       {}
func (g *GameState) SocietyCrumble(pid PID)                    {}
func (g *GameState) SocietyUpdateHealth(pid PID, delta int)    {}
func (g *GameState) SocietyUpdateWealth(pid PID, delta int)    {}
func (g *GameState) SocietyUpdateStability(pid PID, delta int) {}

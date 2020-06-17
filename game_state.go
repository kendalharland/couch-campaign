package couchcampaign

import (
	"couchcampaign/multiplayer"
	"fmt"
)

type gameState struct {
	societies map[multiplayer.CID]SocietyState
	decks     map[multiplayer.CID]*Deck
}

func newGameState(cids []multiplayer.CID) *gameState {
	g := &gameState{
		societies: make(map[multiplayer.CID]SocietyState),
		decks:     make(map[multiplayer.CID]*Deck),
	}
	for _, cid := range cids {
		g.societies[cid] = newSocietyState()
		g.decks[cid] = newDeck()
		stories.Get(storyBasic.ref).AddToDeck(g.decks[cid])
	}
	return g
}

func (g *gameState) HandleInput(cid multiplayer.CID, input Input) (map[multiplayer.CID]SocietyState, error) {
	card := getCard(g.societies[cid].CardRef)

	oldSocieties := make(map[multiplayer.CID]SocietyState)
	for cid := range g.societies {
		oldSocieties[cid] = g.societies[cid]
	}

	switch input {
	case InputCardShown:
		g.applyCardEffects(cid, card.OnShow)
	case InputCardAccepted:
		g.applyCardEffects(cid, card.OnAccept)
	case InputCardRejected:
		g.applyCardEffects(cid, card.OnReject)
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}

	newSocieties := make(map[multiplayer.CID]SocietyState)
	for cid := range oldSocieties {
		if oldSocieties[cid] != g.societies[cid] {
			newSocieties[cid] = g.societies[cid]
		}
	}
	return newSocieties, nil
}

func (g *gameState) applyCardEffects(cid multiplayer.CID, effects []cardEffect) {
	for _, effect := range effects {
		switch e := effect.(type) {
		case AddStoryEffect:
			g.addStory(cid, e)
		case UpdateSocietyStatsEffect:
			g.updateSocietyStats(cid, e)
		default:
			err := fmt.Errorf("unimplemented card effect type %T: %v", e, e)
			panic(err)
		}
	}
}

func (g *gameState) updateSocietyStats(cid multiplayer.CID, e UpdateSocietyStatsEffect) {
	society := g.societies[cid]
	society.Wealth += e.DWealth
	society.Health += e.DHealth
	society.Stability += e.DStability
	g.societies[cid] = society
}

func (g *gameState) addStory(cid multiplayer.CID, e AddStoryEffect) {
	deck := g.decks[cid]
	story := stories.Get(e.StoryRef)
	story.AddToDeck(deck)
}

type governingState struct {
	cid         multiplayer.CID
	actionCount int
}

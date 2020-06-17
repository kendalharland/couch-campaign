package couchcampaign

import (
	"couchcampaign/multiplayer"
	"couchcampaign/proto"
	"fmt"
)

type gameState struct {
	societies   map[multiplayer.CID]SocietyState
	decks       map[multiplayer.CID]*Deck
	isVoting    []multiplayer.CID
	isGoverning []governingState
	updates     []ClientUpdate
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
		g.isGoverning = append(g.isGoverning, governingState{cid, 0})
	}
	return g
}

func (g *gameState) HandleInput(cid multiplayer.CID, input Input) ([]proto.PlayerState, error) {
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

	var newPlayerStates []proto.PlayerState
	for cid := range oldSocieties {
		if oldSocieties[cid] != g.societies[cid] {
			newPlayerStates = append(newPlayerStates, g.getPlayerState(cid))
		}
	}

	return newPlayerStates, nil
}

func (g *gameState) getPlayerState(cid multiplayer.CID) proto.PlayerState {
	society := g.societies[cid]
	card := getCard(society.CardRef)
	return proto.PlayerState{
		Id:                 cid,
		Leader:             society.Leader,
		LeaderTimeInOffice: int32(society.LeaderTimeInOffice),
		Wealth:             int32(society.Wealth),
		Health:             int32(society.Health),
		Stability:          int32(society.Stability),
		Card: &proto.Card{
			Text:       card.Text,
			Speaker:    card.Speaker,
			AcceptText: card.AcceptText,
			RejectText: card.RejectText,
		},
	}
}

func (g *gameState) applyCardEffects(cid multiplayer.CID, effects []cardEffect) {
	for _, effect := range effects {
		switch e := effect.(type) {
		case AddStoryEffect:
			g.addStory(cid, e)
		case SetIsVotingEffect:
			g.setIsVoting(cid)
		case UpdateActionCountEffect:
			g.updateActionCount(cid)
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

func (g *gameState) setIsVoting(cid multiplayer.CID) {
	g.isVoting = append(g.isVoting, cid)
	for i, s := range g.isGoverning {
		if s.cid == cid {
			g.isGoverning = append(g.isGoverning[i:], g.isGoverning[:i+1]...)
		}
	}
}

func (g *gameState) setIsGoverning(cid multiplayer.CID) {
	g.isGoverning = append(g.isGoverning, governingState{cid, 0})
	for i, id := range g.isVoting {
		if id == cid {
			g.isVoting = append(g.isVoting[i:], g.isVoting[:i+1]...)
		}
	}
}

func (g *gameState) updateActionCount(cid multiplayer.CID) {
	for _, s := range g.isGoverning {
		if s.cid == cid {
			s.actionCount++
			return
		}
	}
}

type governingState struct {
	cid         multiplayer.CID
	actionCount int
}

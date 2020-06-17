package couchcampaign

import (
	"fmt"
	"log"

	"couchcampaign/multiplayer"

	"google.golang.org/protobuf/proto"

	couchcampaignpb "couchcampaign/proto"
)

// Game is an instance of a couchcampaign game.
//
// It implements the game's presentation layer.
type game struct {
	state *gameState
}

// NewGame returns a new couchcampaign Game.
func NewGame(state *gameState) multiplayer.Game {
	return &game{state: state}
}

func NewGameWithCIDs(cids []multiplayer.CID) multiplayer.Game {
	state := newGameState(cids)
	return NewGame(state)
}

// HandleMessage handles a client message.
func (g *game) HandleMessage(m multiplayer.Message) ([]multiplayer.Message, error) {
	input, err := parseInput(m.Data)
	if err != nil {
		return nil, fmt.Errorf("parseInput: %w", err)
	}

	societies, err := g.state.HandleInput(m.CID, input)
	if err != nil {
		return nil, fmt.Errorf("HandleUpdate: %w", err)
	}

	messages := make([]multiplayer.Message, 0, len(societies))
	for cid, s := range societies {
		ps := societyToPlayerState(cid, s)
		psData, err := proto.Marshal(&ps)
		if err != nil {
			return nil, err
		}
		messages = append(messages, multiplayer.Message{
			CID:  ps.Id,
			Data: psData,
		})
	}
	return messages, nil
}

// HandleError handles a client error.
func (g *game) HandleError(e multiplayer.ClientError) error {
	log.Printf("error: %v: %v", e.CID, e.Err)
	return nil
}

// Close disposes this game and all of its resources.
func (g *game) Close() {}

func parseInput(input []byte) (Input, error) {
	value := string(input)
	switch value {
	case "accept":
		return InputCardAccepted, nil
	case "reject":
		return InputCardRejected, nil
	case "show":
		return InputCardShown, nil
	default:
		return InputErr, fmt.Errorf("invalid input: %q", value)
	}
}

func societyToPlayerState(cid multiplayer.CID, society SocietyState) couchcampaignpb.PlayerState {
	card := getCard(society.CardRef)
	return couchcampaignpb.PlayerState{
		Id:                 cid,
		Leader:             society.Leader,
		LeaderTimeInOffice: int32(society.LeaderTimeInOffice),
		Wealth:             int32(society.Wealth),
		Health:             int32(society.Health),
		Stability:          int32(society.Stability),
		Card: &couchcampaignpb.Card{
			Text:       card.Text,
			Speaker:    card.Speaker,
			AcceptText: card.AcceptText,
			RejectText: card.RejectText,
		},
	}
}

package couchcampaign

import (
	"couchcampaign/multiplayer"
	"fmt"
	"log"
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
func (g *game) HandleMessage(m multiplayer.Message) error {
	input, err := parseInput(m.Data)
	if err != nil {
		return fmt.Errorf("parseInput: %w", err)
	}
	return g.state.HandleInput(m.CID, input)
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

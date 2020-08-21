package couchcampaign

import (
	"fmt"

	"couchcampaign/starlarkgame"
)

const scriptsFilename = "scripts/main.star"

type Game struct {
	g *starlarkgame.Game
}

func NewGame() (*Game, error) {
	ctx := starlarkgame.NewContext()
	g, err := starlarkgame.New(ctx, scriptsFilename)
	if err != nil {
		return nil, err
	}
	return &Game{g}, nil
}

func (g *Game) HandleInput(data []byte) error {
	input, err := parseInput(data)
	if err != nil {
		return fmt.Errorf("parseInput: %w", err)
	}
	if err := g.g.HandleInput(string(input)); err != nil {
		return fmt.Errorf("HandleInput(%v): %w", input, err)
	}
	return nil
}

func (g *Game) GetPlayerState() starlarkgame.PlayerState {
	return g.g.GetPlayerState()
}

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

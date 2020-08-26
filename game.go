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
	ctx := starlarkgame.NewContext(*starlarkgame.NewDeck())
	g, err := starlarkgame.New(ctx, scriptsFilename)
	if err != nil {
		return nil, err
	}
	return &Game{g}, nil
}

func (g *Game) HandleInput(input Input) error {
	if err := g.g.HandleInput(string(input)); err != nil {
		return fmt.Errorf("HandleInput(%v): %w", input, err)
	}
	return nil
}

func (g *Game) GetPlayerState() starlarkgame.PlayerState {
	return g.g.GetPlayerState()
}

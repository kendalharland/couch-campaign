package starlarkgame

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

const (
	fnAddPlayer   = "add_player"
	fnHandleInput = "handle_input"
	fnNewGame     = "new_game"
	fnLoadCard    = "load_card"
)

type Game struct {
	interpreter *Interpreter
	state       starlark.Value
	ctx         *Context
}

func New(ctx *Context, filename string) (*Game, error) {
	i := NewInterpreter(builtins(ctx))
	if err := i.ExecFile(filename); err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", filename, err)
	}
	state, err := i.Call(fnNewGame, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new game: %w", err)
	}
	return &Game{ctx: ctx, interpreter: i, state: state}, nil
}

func (g *Game) GetPlayerState() PlayerState {
	return *(g.ctx.GetPlayerState())
}

func (g *Game) GetCurrentCard() (*Card, error) {
	// topCard := string(g.ctx.TopCard())
	// args := starlark.Tuple{g.state, starlark.String(topCard)}
	// value, err := g.interpreter.Call(fnLoadCard, args, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// return value.(*Card), nil
	return nil, errors.New("unimplemented: GetCurrentCard")
}

func (g *Game) HandleInput(input string) error {
	args := starlark.Tuple{g.state, starlark.String(input)}
	state, err := g.interpreter.Call(fnHandleInput, args, nil)
	if err != nil {
		return err
	}
	g.state = state
	return nil
}

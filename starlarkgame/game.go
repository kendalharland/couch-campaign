package starlarkgame

import (
	"fmt"

	"go.starlark.net/starlark"
)

const (
	fnAddPlayer   = "add_player"
	fnHandleInput = "handle_input"
	fnNewGame     = "new_game"
)

type Game struct {
	interpreter *Interpreter
	state       starlark.Value
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
	return &Game{interpreter: i, state: state}, nil
}

func (g *Game) HandleInput(playerID, input string) error {
	args := starlark.Tuple{g.state, starlark.String(playerID), starlark.String(input)}
	state, err := g.interpreter.Call(fnHandleInput, args, nil)
	if err != nil {
		return err
	}
	g.state = state
	return nil
}

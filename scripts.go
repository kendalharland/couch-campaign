package couchcampaign

import (
	"couchcampaign/starlarkgame"

	"go.starlark.net/starlark"
)

const (
	fnAddPlayer   = "add_player"
	fnHandleInput = "handle_input"
	fnNewGame     = "new_game"
)

type scriptManager struct {
	interpreter *starlarkgame.Interpreter

	// The Game state as seen by the starlark scripts.
	//
	// The script manager does not modify this, it only passes it between invocations of
	// the game scripts.
	gameState starlark.Value
}

func newScriptManager(ctx *gameContext) *scriptManager {
	return &scriptManager{
		interpreter: starlarkgame.NewInterpreter(builtins(ctx)),
	}
}

func (g *scriptManager) LoadMainScript(filename string) error {
	if err := g.interpreter.ExecFile(filename); err != nil {
		return err
	}
	gameState, err := g.interpreter.Call(fnNewGame, nil, nil)
	if err != nil {
		return err
	}
	g.gameState = gameState
	return nil
}

func (g *scriptManager) HandleInput(playerID string, input Input) error {
	_, err := g.interpreter.Call(fnHandleInput, starlark.Tuple{g.gameState, starlark.String(playerID), starlark.String(input)}, nil)
	return err
}

func builtins(ctx *gameContext) starlark.StringDict {
	return starlark.StringDict{
		"_builtins_get_player_ids":       starlark.NewBuiltin("_builtins_get_player_ids", _builtins_get_player_ids(ctx)),
		"_builtins_get_player_card_id":   starlark.NewBuiltin("_builtins_get_player_card_id", _builtins_get_player_card_id(ctx)),
		"_builtins_set_player_card_id":   starlark.NewBuiltin("_builtins_set_player_card_id", _builtins_set_player_card_id(ctx)),
		"_builtins_add_player_health":    starlark.NewBuiltin("_builtins_add_player_health", _builtins_add_player_health(ctx)),
		"_builtins_add_player_wealth":    starlark.NewBuiltin("_builtins_add_player_wealth", _builtins_add_player_wealth(ctx)),
		"_builtins_add_player_stability": starlark.NewBuiltin("_builtins_add_player_stability", _builtins_add_player_stability(ctx)),
	}
}

type starlarkBuiltin = func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)

func _builtins_get_player_ids(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}
		playerIDs := ctx.PlayerIDs()
		values := make([]starlark.Value, 0, len(playerIDs))
		for _, id := range playerIDs {
			values = append(values, starlark.String(id))
		}
		return starlark.NewList(values), nil
	}
}

func _builtins_get_player_card_id(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var playerID string
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &playerID); err != nil {
			return nil, err
		}
		playerState, err := ctx.GetPlayer(playerID)
		if err != nil {
			return nil, err
		}
		return starlark.String(playerState.CardID), nil
	}
}

func _builtins_set_player_card_id(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var playerID, cardID string
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 2, &playerID, &cardID); err != nil {
			return nil, err
		}
		playerState, err := ctx.GetPlayer(playerID)
		if err != nil {
			return nil, err
		}
		playerState.CardID = cardID
		if err := ctx.SetPlayer(playerID, *playerState); err != nil {
			return nil, err
		}
		return starlark.None, nil
	}
}

func _builtins_add_player_health(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var playerID string
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 2, &playerID, &amount); err != nil {
			return nil, err
		}
		playerState, err := ctx.GetPlayer(playerID)
		if err != nil {
			return nil, err
		}
		playerState.Health += amount
		if err := ctx.SetPlayer(playerID, *playerState); err != nil {
			return nil, err
		}
		return starlark.None, nil
	}
}

func _builtins_add_player_wealth(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var playerID string
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 2, &playerID, &amount); err != nil {
			return nil, err
		}
		playerState, err := ctx.GetPlayer(playerID)
		if err != nil {
			return nil, err
		}
		playerState.Wealth += amount
		if err := ctx.SetPlayer(playerID, *playerState); err != nil {
			return nil, err
		}
		return starlark.None, nil
	}
}

func _builtins_add_player_stability(ctx *gameContext) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var playerID string
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 2, &playerID, &amount); err != nil {
			return nil, err
		}
		playerState, err := ctx.GetPlayer(playerID)
		if err != nil {
			return nil, err
		}
		playerState.Stability += amount
		if err := ctx.SetPlayer(playerID, *playerState); err != nil {
			return nil, err
		}
		return starlark.None, nil
	}
}

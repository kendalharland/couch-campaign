package starlarkgame

import "go.starlark.net/starlark"

func builtins(ctx *Context) starlark.StringDict {
	return starlark.StringDict{
		"_builtins_get_player_card_id":   starlark.NewBuiltin("_builtins_get_player_card_id", _builtins_get_player_card_id(ctx)),
		"_builtins_set_player_card_id":   starlark.NewBuiltin("_builtins_set_player_card_id", _builtins_set_player_card_id(ctx)),
		"_builtins_add_player_health":    starlark.NewBuiltin("_builtins_add_player_health", _builtins_add_player_health(ctx)),
		"_builtins_add_player_wealth":    starlark.NewBuiltin("_builtins_add_player_wealth", _builtins_add_player_wealth(ctx)),
		"_builtins_add_player_stability": starlark.NewBuiltin("_builtins_add_player_stability", _builtins_add_player_stability(ctx)),
	}
}

type starlarkBuiltin = func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)

func _builtins_get_player_card_id(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}
		return starlark.String(ctx.GetPlayerState().CardID), nil
	}
}

func _builtins_set_player_card_id(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var cardID string
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &cardID); err != nil {
			return nil, err
		}
		ps := ctx.GetPlayerState()
		ps.CardID = cardID
		ctx.SetPlayerState(ps)
		return starlark.None, nil
	}
}

func _builtins_add_player_health(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &amount); err != nil {
			return nil, err
		}
		ps := ctx.GetPlayerState()
		ps.Health += amount
		ctx.SetPlayerState(ps)
		return starlark.None, nil
	}
}

func _builtins_add_player_wealth(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &amount); err != nil {
			return nil, err
		}
		ps := ctx.GetPlayerState()
		ps.Wealth += amount
		ctx.SetPlayerState(ps)
		return starlark.None, nil
	}
}

func _builtins_add_player_stability(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var amount int
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &amount); err != nil {
			return nil, err
		}
		ps := ctx.GetPlayerState()
		ps.Stability += amount
		ctx.SetPlayerState(ps)
		return starlark.None, nil
	}
}

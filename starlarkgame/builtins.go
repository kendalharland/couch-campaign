package starlarkgame

import (
	"errors"

	"go.starlark.net/starlark"
)

//
func builtins(ctx *Context) starlark.StringDict {
	return starlark.StringDict{
		"_builtins_get_player":     starlark.NewBuiltin("_builtins_get_player", _builtins_get_player(ctx)),
		"_builtins_deck_push_card": starlark.NewBuiltin("_builtins_deck_push_card", _builtins_deck_push_card(ctx)),
		"_builtins_deck_pop_card":  starlark.NewBuiltin("_builtins_deck_pop_card", _builtins_deck_pop_card(ctx)),
	}
}

type starlarkBuiltin = func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error)

func _builtins_deck_push_card(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var cardRef string
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 1, &cardRef); err != nil {
			return nil, err
		}
		if cardRef == "" {
			return nil, errors.New("card ref cannot be the empty string")
		}
		ctx.DeckPushCard(cardRef)
		return starlark.None, nil
	}
}

func _builtins_deck_pop_card(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}
		ctx.DeckPopCard()
		return starlark.None, nil
	}
}

func _builtins_get_player(ctx *Context) starlarkBuiltin {
	return func(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		if err := starlark.UnpackPositionalArgs(f.Name(), args, kwargs, 0); err != nil {
			return nil, err
		}
		ps := ctx.GetPlayerState()
		return newGoStruct("player_state", ps), nil
	}
}

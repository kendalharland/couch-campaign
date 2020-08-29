package couchcampaign

import (
	"couchcampaign/starlarkgame"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGame(t *testing.T) {
	game, err := NewGame()
	if err != nil {
		t.Fatalf("NewGame got unexpected error: %v", err)
	}

	s := game.GetState()
	want := starlarkgame.State{
		CardRef:            "viral_infection",
		Leader:             "",
		LeaderTimeInOffice: 0,
		Health:             0,
		Wealth:             0,
		Stability:          0,
	}
	if diff := cmp.Diff(s, want); diff != "" {
		t.Fatalf("got diff: (+got,-want)\n%s\n", diff)
	}

	if err := game.HandleInput(InputCardAccepted); err != nil {
		t.Fatalf("HandleInput(%q) unexpected error %v", InputCardAccepted, err)
	}

	s = game.GetState()
	want = starlarkgame.State{
		CardRef:            "tobbacco_ad",
		Leader:             "",
		LeaderTimeInOffice: 0,
		Health:             -2,
		Wealth:             2,
		Stability:          0,
	}
	if diff := cmp.Diff(s, want); diff != "" {
		t.Fatalf("got diff: (+got,-want)\n%s\n", diff)
	}

	if err := game.HandleInput(InputCardRejected); err != nil {
		t.Fatalf("HandleInput(%q) unexpected error %v", InputCardRejected, err)
	}

	s = game.GetState()
	want = starlarkgame.State{
		CardRef:            "",
		Leader:             "",
		LeaderTimeInOffice: 0,
		Health:             -4,
		Wealth:             4,
		Stability:          0,
	}
	if diff := cmp.Diff(s, want); diff != "" {
		t.Fatalf("got diff: (+got,-want)\n%s\n", diff)
	}
}

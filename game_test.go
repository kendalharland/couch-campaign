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

	ps := game.GetPlayerState()
	want := starlarkgame.PlayerState{
		// Card attributes.
		CardID:         "viral_infection",
		CardSpeaker:    "",
		CardText:       "",
		CardAcceptText: "",
		CardRejectText: "",
		// Societal attributes.
		Leader:             "",
		LeaderTimeInOffice: 0,
		Health:             0,
		Wealth:             0,
		Stability:          0,
	}
	if diff := cmp.Diff(ps, want); diff != "" {
		t.Fatalf("got diff: (-want,+got)\n%s\n", diff)
	}

	if err := game.HandleInput([]byte(InputCardAccepted)); err != nil {
		t.Fatalf("HandleInput(%q) unexpected error %v", InputCardAccepted, err)
	}

	// TODO: The card should have been updated from the deck here.
	ps = game.GetPlayerState()
	want = starlarkgame.PlayerState{
		// Card attributes.
		CardID:         "viral_infection",
		CardSpeaker:    "",
		CardText:       "",
		CardAcceptText: "",
		CardRejectText: "",
		// Societal attributes.
		Leader:             "",
		LeaderTimeInOffice: 0,
		Health:             -2,
		Wealth:             2,
		Stability:          0,
	}
	if diff := cmp.Diff(ps, want); diff != "" {
		t.Fatalf("got diff: (-want,+got)\n%s\n", diff)
	}
}

package couchcampaign

import (
	"testing"
)

const (
	p1 = "p1"
	p2 = "p2"
)

func TestGame(t *testing.T) {
	game, err := NewGame([]string{p1, p2})
	if err != nil {
		t.Fatalf("NewGame got unexpected error: %v", err)
	}
	if err := game.HandleInput(p1, []byte(InputCardAccepted)); err != nil {
		t.Fatalf("HandleInput(%q, %q) unexpected error %v", p1, InputCardAccepted, err)
	}
	ps, err := game.GetPlayerState(p1)
	if err != nil {
		t.Fatalf("GetPlayerState(%q) unexpected error: %v", ps, err)
	}
}

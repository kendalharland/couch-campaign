package couchcampaign_test

import (
	"couchcampaign"
	"couchcampaign/multiplayer"
	"testing"
)

func TestGame(t *testing.T) {
	game := couchcampaign.NewGame()
	game.AddPlayer("p1")

	if err := game.Start(); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	m := multiplayer.Message{CID: "p1", Data: []byte("accept")}
	s, err := game.HandleMessage(m)
	if err != nil {
		t.Fatalf("failed to handle input %q: %v", m.Data, err)
	}

	if len(s) == 0 {
		t.Errorf("wanted player state change but got %v", s)
	}
}

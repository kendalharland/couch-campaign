package couchcampaign_test

import (
	"couchcampaign"
	"testing"
)

func TestGame(t *testing.T) {
	game := couchcampaign.NewGame()
	game.AddPlayer("p1")

	if err := game.Start(); err != nil {
		t.Fatalf("failed to start game: %v", err)
	}

	input := "accept"

	if err := game.HandleInput("p1", []byte(input)); err != nil {
		t.Fatalf("failed to handle input %q: %v", input, err)
	}

	message := <-game.Outputs()
	outputCID := message.CID
	outputData := message.Data

	if outputCID != "p1" || len(outputData) == 0 {
		t.Errorf("unexpected message sent (cid=%q, data=%q)", outputCID, string(outputData))
	}
}

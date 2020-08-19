package starlarkgame

import "testing"

func TestGame(t *testing.T) {
	gtx := NewContext()
	gtx.AddPlayer("player1")

	game, err := New(gtx, "test_scripts/test_game.star")
	if err != nil {
		t.Fatalf("test setup failed: %v", err)
	}

	if err := game.HandleInput("player1", ""); err != nil {
		t.Fatalf("HandleInput(player1, '') got unexpected error: %v", err)
	}
	if err := game.HandleInput("nosuchplayer", ""); err == nil {
		t.Fatal("HandleInput(nosuchplayer, '') wanted an error but got nil")
	}
}

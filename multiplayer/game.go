package multiplayer

// Game implements the core game logic.
type Game interface {
	HandleInput(CID, []byte) error
	GetPlayerState(CID) ([]byte, error)
}

package multiplayer

// Game implements the core game logic.
type Game interface {
	AddPlayer(CID) error
	Start() error
	Stop() error
	HandleInput(CID, []byte) error
	HandleError(ClientError) error
	Outputs() <-chan Message
}

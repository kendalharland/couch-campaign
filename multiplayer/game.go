package multiplayer

// Game implements the core game logic.
type Game interface {
	HandleMessage(Message) ([]Message, error)
	HandleError(ClientError) error
	Close()
}

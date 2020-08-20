package multiplayer

type Game interface {
	HandleInput(CID, []byte) error
	GetPlayerState(CID) ([]byte, error)
}

type GameBuilder func([]CID) (Game, error)

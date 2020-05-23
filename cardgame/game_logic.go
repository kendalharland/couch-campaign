package cardgame

type GameLogic interface {
	OnCardAccepted(GameController, PID, CardRef) ([]PID, error)
	OnCardRejected(GameController, PID, CardRef) ([]PID, error)
	OnCardShown(GameController, PID, CardRef) ([]PID, error)
}

type GameController interface {
	DeckInsertAll(c CardRef)
	DeckInsert(pid PID, card string)
	DeckRemove(pid PID, card string)
	DeckShuffle(pid PID)
	SocietyCrumble(pid PID)
	SocietyUpdateHealth(pid PID, delta int)
	SocietyUpdateWealth(pid PID, delta int)
	SocietyUpdateStability(pid PID, delta int)
}

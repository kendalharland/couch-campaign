package couchcampaign

type GameLogicApi interface {
	DeckInsertAll(c CardRef)
	DeckInsert(pid PID, card string)
	DeckRemove(pid PID, card string)
	DeckShuffle(pid PID)
	SocietyCrumble(pid PID)
	SocietyUpdateHealth(pid PID, delta int)
	SocietyUpdateWealth(pid PID, delta int)
	SocietyUpdateStability(pid PID, delta int)
}

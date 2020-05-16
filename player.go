package couchcampaign

import (
	"log"
	"math/rand"

	"github.com/gobuffalo/uuid"
	"github.com/gonum/stat"
)

// Leader names
var leaders = []string{
	"Argal Essential",
	"Ate Predicate",
	"Bashful Bailey",
	"Beef Beetfly",
	"Boogie Boost",
	"Boogie Bothy",
	"Boogie Bowler",
	"Capital Vital",
	"Chide Fried",
	"Chorale Chough",
	"Cooking Cooked",
	"Crucial Credo",
	"Dithyramb Lamb",
	"Eats Eave",
	"Eats Receipts",
	"Essential Essay",
	"Harp Hanky",
	"Harpy Harmonica",
	"Hesitant Heater",
	"Hesitant Heinie",
	"Key Precis",
	"Lamb Fram",
	"Leben Eaten",
	"Leery Leben",
	"Licks Ligan",
	"Loment Reticent",
	"Percent Reticent",
	"Pork Bork",
	"Ragtime Climb",
	"Rate Ate",
}

// PID is a type alias for a player ID.
type PID string

// NewPID generates a unique PID.
func NewPID() PID {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("newPID: %v", err)
	}
	return PID(id.String())
}

func (p PID) String() string {
	return string(p)
}

const (
	initStatValue = 15
	maxStatValue  = 30
	minStatValue  = 0
)

// EmptyPlayerState is returned when a player's society has collapsed.
var EmptyPlayerState = playerState{}

type playerState struct {
	Card      Card
	Leader    string
	Wealth    int
	Stability int
	Health    int
}

func newPlayerState() playerState {
	return playerState{
		Card:      nil,
		Leader:    leaders[rand.Intn(len(leaders))],
		Wealth:    initStatValue,
		Health:    initStatValue,
		Stability: initStatValue,
	}
}

func OnVotingCard(s playerState, _ votingCard) (next playerState) {
	next = s
	return checkPlayerState(next)
}

func OnDismissInfoCard(s playerState, _ infoCard) (next playerState) {
	next = s
	return checkPlayerState(next)
}

func OnAcceptActionCard(s playerState, c actionCard) (next playerState) {
	next.Wealth = s.Wealth + c.AccWealthEffect
	next.Health = s.Health + c.AccHealthEffect
	next.Stability = s.Stability + c.AccStabilityEffect
	return checkPlayerState(next)
}

func OnRejectActionCard(s playerState, c actionCard) (next playerState) {
	next.Wealth = s.Wealth + c.RejWealthEffect
	next.Health = s.Health + c.RejHealthEffect
	next.Stability = s.Stability + c.RejStabilityEffect
	return checkPlayerState(next)
}

func checkPlayerState(s playerState) playerState {
	switch {
	case s.Wealth <= minStatValue || maxStatValue <= s.Wealth:
		return EmptyPlayerState
	case s.Health <= minStatValue || maxStatValue <= s.Health:
		return EmptyPlayerState
	case s.Stability <= minStatValue || maxStatValue <= s.Stability:
		return EmptyPlayerState
	default:
		return s
	}
}

func (s playerState) SocietyScore() int {
	return s.Wealth + s.Stability + s.Health
}

func (s playerState) SocietyVariance() float64 {
	return stat.Variance([]float64{
		float64(s.Wealth),
		float64(s.Stability),
		float64(s.Health),
	}, nil)
}

package cardgame

import (
	"log"
	"math/rand"

	"github.com/gobuffalo/uuid"
	"github.com/gonum/stat"
)

// Leader names
//
// TODO This should be part of the client library.
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

// NilSocietyState is returned when a player's society has collapsed.
var NilSocietyState = SocietyState{}

// SocietyState is only modified when the client that owns this state sends
// a message to the server.
type SocietyState struct {
	CardRef   CardRef
	Leader    string
	Wealth    int
	Stability int
	Health    int
}

func newSocietyState() SocietyState {
	return SocietyState{
		CardRef:   "",
		Leader:    leaders[rand.Intn(len(leaders))],
		Wealth:    initStatValue,
		Health:    initStatValue,
		Stability: initStatValue,
	}
}

func checkSocietyState(s SocietyState) SocietyState {
	switch {
	case s.Wealth <= minStatValue || maxStatValue <= s.Wealth:
		return NilSocietyState
	case s.Health <= minStatValue || maxStatValue <= s.Health:
		return NilSocietyState
	case s.Stability <= minStatValue || maxStatValue <= s.Stability:
		return NilSocietyState
	default:
		return s
	}
}

func (s SocietyState) SocietyScore() int {
	return s.Wealth + s.Stability + s.Health
}

func (s SocietyState) SocietyVariance() float64 {
	return stat.Variance([]float64{
		float64(s.Wealth),
		float64(s.Stability),
		float64(s.Health),
	}, nil)
}

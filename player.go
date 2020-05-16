package couchcampaign

import (
	"log"
	"math/rand"

	"github.com/gobuffalo/uuid"
	"github.com/gonum/stat"
)

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

const (
	initStatValue = 15
	maxStatValue  = 30
	minStatValue  = 0
)

type stats struct {
	Leader    string
	Wealth    int
	Stability int
	Health    int
}

func newStats() *stats {
	return &stats{
		Leader:    leaders[rand.Intn(len(leaders))],
		Wealth:    initStatValue,
		Health:    initStatValue,
		Stability: initStatValue,
	}
}

func (s stats) Sum() int {
	return s.Wealth + s.Stability + s.Health
}

func (s stats) Variance() float64 {
	return stat.Variance([]float64{
		float64(s.Wealth),
		float64(s.Stability),
		float64(s.Health),
	}, nil)
}

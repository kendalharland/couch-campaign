package couchcampaign

// Leader names.
//
// One of these is selected at random after the player's society inevitably crumbles and
// they are disgracefully thrown out of office.
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

type SocietyState struct {
	CardRef            CardRef
	Leader             string
	LeaderTimeInOffice int
	Wealth             int
	Stability          int
	Health             int
}

func newSocietyState() SocietyState {
	return SocietyState{
		CardRef:   "",
		Leader:    leaders[0],
		Wealth:    initStatValue,
		Health:    initStatValue,
		Stability: initStatValue,
	}
}

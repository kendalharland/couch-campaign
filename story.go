package couchcampaign

// Global story data.
//
// These are initialized at startup and shoul not be modified while the game is
// running.
var (
	stories = newStoryRegistry()
	cards   = map[CardRef]Card{
		// Core cards.
		"WaitingForVotes": {
			ID:   "WaitingForVotes",
			Text: "Waiting for votes...",
			OnShow: []cardEffect{
				SetIsVotingEffect{},
			},
		},
	}
)

func getCard(c CardRef) Card {
	for _, card := range cards {
		if card.ID == c {
			return card
		}
	}
	return Card{}
}

type storyRef = string

type story struct {
	ref   string
	cards []Card
}

func (s story) AddToDeck(d *Deck) {
	for _, card := range s.cards {
		d.Insert(CardNode{
			Card:     card.ID,
			Priority: MinCardPriority,
		})
	}
}

type storyRegistry struct {
	stories map[storyRef]story
}

func newStoryRegistry() *storyRegistry {
	return &storyRegistry{
		stories: make(map[storyRef]story),
	}
}

func (r *storyRegistry) Register(s story) {
	r.stories[s.ref] = s
}

func (r *storyRegistry) Get(s storyRef) story {
	return r.stories[s]
}

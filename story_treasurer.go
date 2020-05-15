package couchcampaign

var storyTreasurer = story{
	ref: "treasurer",
	cards: []Card{{
		ID:         "RaiseTaxes",
		Speaker:    "Treasurer",
		Text:       "We need to impose new taxes on the rich to pay for new schools in the area",
		AcceptText: "Make it happen",
		OnAccept: []cardEffect{
			UpdateSocietyStatsEffect{3, 0, -1},
			UpdateActionCountEffect{},
		},
		RejectText: "My friends would not appreciate that",
		OnReject: []cardEffect{
			UpdateSocietyStatsEffect{-1, 0, -2},
			UpdateActionCountEffect{},
		},
	}},
}

func init() {
	stories.Register(storyTreasurer)
}

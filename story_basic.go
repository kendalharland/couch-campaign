package couchcampaign

var storyBasic = story{
	ref: "basic",
	cards: []Card{
		// Basic event cards.
		{
			ID:         "ViralInfection",
			Speaker:    "SurgeonGeneral",
			Text:       "A novel viral infection broke out among a cargo ship waiting to come ashore. What should we do?",
			AcceptText: "Quarantine the ship in the harbor.",
			OnAccept: []cardEffect{
				UpdateSocietyStatsEffect{-2, 2, 0},
				UpdateActionCountEffect{},
			},
			RejectText: "Have them deliver their goods, then see them to the nearest hospital immediately.",
			OnReject: []cardEffect{
				UpdateSocietyStatsEffect{-3, -2, 0},
				UpdateActionCountEffect{},
			},
		},
		{
			ID:         "TobaccoAd",
			Speaker:    "SurgeonGeneral",
			Text:       "A tobacco company wants to advertise at the community center. They're offering to cut us in on the profits...",
			AcceptText: "Our coffers have been running a little dry...",
			OnAccept: []cardEffect{
				UpdateSocietyStatsEffect{3, -2, -1},
				UpdateActionCountEffect{},
			},
			RejectText: "I won't sacrifice the public health for financial gain.",
			OnReject: []cardEffect{
				UpdateSocietyStatsEffect{1, 1, 0},
				UpdateActionCountEffect{},
			},
		},
		// Recruiting cards.
		{
			ID:         "RecruitTreasurer",
			Speaker:    "Treasurer",
			Text:       "Your office is bankrupt and you are terrible at managing your finances. Let me help",
			AcceptText: "Erm... Sure",
			OnAccept: []cardEffect{
				UpdateSocietyStatsEffect{0, 0, 2},
				AddStoryEffect{storyTreasurer.ref},
				UpdateActionCountEffect{},
			},
			RejectText: "I can spend the people's money however I see fit!",
			OnReject: []cardEffect{
				UpdateActionCountEffect{},
			},
		},
	},
}

func init() {
	stories.Register(storyBasic)
}

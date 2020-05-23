package couchcampaign

import "couchcampaign/cardgame"

const (
	self                    = "[self]"
	speakerCommunityPlanner = "CommunityPlanner"
	speakerDistrictAttorney = "District Attorney"
	speakerSurgeonGeneral   = "Surgeon General"
)

var (
	TheVotingCard = cardgame.Card{
		ID:   "TheVotingCard",
		Text: "waiting for votes...",
		OnShown: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			// api.ElectionIsVoting(pid)
		},
	}

	TheWelcomeCard = cardgame.Card{
		ID: "TheWelcomeCard",
		Text: "Welcome to couchcampaign!\n\n" +
			"Today is your first day in office as the governor.\n" +
			"Work with your advisors to stay in office as long as you can.\n" +
			"Remember, the ultimate goal is to win the presidency.\n.",
	}

	SocietyCrumbleCard = cardgame.Card{
		ID:   "SocietyCrumbleCard",
		Text: "Society has crumbled and you are being forced out of office.",
	}

	ElectionCampaignSeasonCard = cardgame.Card{
		ID:   "ElectionCampaignSeasonCard",
		Text: "Campaign season has begun!",
	}

	ElectionVotingSeasonCard = cardgame.Card{
		ID:   "ElectionVotingSeasonCard",
		Text: "Voting season has begun!",
	}

	ElectionOffSeasonCard = cardgame.Card{
		ID:   "ElectionOffSeasonCard",
		Text: "The offseason has begun!",
	}

	ElectionWonCard = cardgame.Card{
		ID:   "ElectionWonCard",
		Text: "You won the election! Now get back to work.",
	}

	ElectionLostCard = cardgame.Card{
		ID:   "ElectionLostCard",
		Text: "You lost the election. Now get back to to work.",
	}

	CommunityPlanner_FestivalFundraiser = cardgame.Card{
		ID:         "CommuityPlanner_FestivalFundraiser",
		Speaker:    speakerCommunityPlanner,
		Text:       "We need some money to organize the annual spring festival this year, it will be amazing",
		AcceptText: "Yes",
		OnAccept: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateWealth(pid, -2)
			api.SocietyUpdateHealth(pid, 1)
			api.SocietyUpdateStability(pid, 0)
			// api.ElectionUpdate(pid)
		},
		RejectText: "No",
		OnReject: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateHealth(self, -2)
			api.SocietyUpdateStability(self, 0)
			// api.ElectionUpdate(pid)
		},
	}

	DistrictAttorney_JailBreak = cardgame.Card{
		ID:         "DistrictAttorney_JailBreak",
		Speaker:    speakerDistrictAttorney,
		Text:       "Some prisoners recently escaped from the local jail. Should we put local precints on alert?",
		AcceptText: "Yes, alert all forces immediately",
		OnAccept: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateHealth(self, 1)
			api.SocietyUpdateStability(self, 2)
			// api.ElectionUpdate(pid)
		},
		RejectText: "No, assign a small task force. We can't let the press find out",
		OnReject: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateHealth(self, -2)
			api.SocietyUpdateStability(self, -2)
			// api.ElectionUpdate(pid)
		},
	}

	SurgeonGeneral_TobaccoAd = cardgame.Card{
		ID:         "SurgeonGeneral_TobaccoAd",
		Speaker:    speakerSurgeonGeneral,
		Text:       "A tobacco company wants to advertise around the community center. They're offering to cut us in on the profits...",
		AcceptText: "Our coffers have been running a little dry...",
		OnAccept: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateWealth(self, 3)
			api.SocietyUpdateHealth(self, -2)
			api.SocietyUpdateStability(self, -1)
			// api.ElectionUpdate(pid)
		},
		RejectText: "I won't sacrifice the public health for financial gain.",
		OnReject: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateHealth(self, 1)
			api.SocietyUpdateStability(self, 1)
			// api.ElectionUpdate(pid)
		},
	}

	SurgeonGeneral_ViralInfection = cardgame.Card{
		ID:         "SurgeonGeneral_ViralInfection",
		Speaker:    speakerSurgeonGeneral,
		Text:       "A novel viral infection broke out among a cargo ship waiting to come ashore. What should we do?",
		AcceptText: "Quarantine the ship in the harbor.",
		OnAccept: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateWealth(self, -2)
			api.SocietyUpdateHealth(self, 2)
			// api.ElectionUpdate(pid)
		},
		RejectText: "Have them deliver their goods, then see them to the nearest hospital immediately.",
		OnReject: func(api cardgame.GameLogicApi, pid cardgame.PID, card cardgame.CardRef) {
			api.SocietyUpdateWealth(self, 3)
			api.SocietyUpdateHealth(self, -2)
			// api.ElectionUpdate(pid)
		},
	}
)

var allCards = map[cardgame.CardRef]cardgame.Card{
	TheVotingCard.ID:                       TheVotingCard,
	TheWelcomeCard.ID:                      TheWelcomeCard,
	ElectionCampaignSeasonCard.ID:          ElectionCampaignSeasonCard,
	ElectionVotingSeasonCard.ID:            ElectionVotingSeasonCard,
	ElectionOffSeasonCard.ID:               ElectionOffSeasonCard,
	ElectionWonCard.ID:                     ElectionWonCard,
	ElectionLostCard.ID:                    ElectionLostCard,
	SocietyCrumbleCard.ID:                  SocietyCrumbleCard,
	CommunityPlanner_FestivalFundraiser.ID: CommunityPlanner_FestivalFundraiser,
	DistrictAttorney_JailBreak.ID:          DistrictAttorney_JailBreak,
	SurgeonGeneral_TobaccoAd.ID:            SurgeonGeneral_TobaccoAd,
	SurgeonGeneral_ViralInfection.ID:       SurgeonGeneral_ViralInfection,
}

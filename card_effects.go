package couchcampaign

type cardEffect interface{}

type AddStoryEffect struct {
	StoryRef storyRef
}

type SetIsVotingEffect struct{}

type UpdateActionCountEffect struct{}

type UpdateSocietyStatsEffect struct {
	DWealth, DHealth, DStability int
}

package couchcampaign

type cardEffect interface{}

type AddStoryEffect struct {
	StoryRef storyRef
}

type UpdateSocietyStatsEffect struct {
	DWealth, DHealth, DStability int
}

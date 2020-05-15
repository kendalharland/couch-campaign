package couchcampaign

import "strings"

var storyIntro = story{
	ref: "intro",
	cards: []Card{{
		ID: "Welcome",
		Text: strings.Join([]string{
			"Welcome to couchcampaign!\n",
			"Today is your first day in office.",
			"Work with your advisors to stay in office as long as you can.",
			"The ultimate goal is to win the presidency.",
		}, "\n"),
	}},
}

func init() {
	stories.Register(storyIntro)
}

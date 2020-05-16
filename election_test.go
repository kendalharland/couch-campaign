package couchcampaign

import (
	"testing"
)

func TestElectionTracker(t *testing.T) {
	p1 := NewPID()
	p2 := NewPID()

	t.Run("2p", func(t *testing.T) {
		type testStep struct {
			pid        PID
			card       Card
			wantSeason electionSeason
		}
		tests := []struct {
			name         string
			offSeasonLen int
			steps        []testStep
		}{
			{
				name:         "not enough cards to trigger campaign season",
				offSeasonLen: 2,
				steps: []testStep{
					{p1, actionCard{}, offSeason},
					{p2, actionCard{}, offSeason},
				},
			},
			{
				name:         "campaign season triggered",
				offSeasonLen: 2,
				steps: []testStep{
					{p1, actionCard{}, offSeason},
					{p2, actionCard{}, offSeason},
					{p1, actionCard{}, campaignSeason},
					{p1, actionCard{}, campaignSeason},
				},
			},
			{
				name:         "voting season triggered",
				offSeasonLen: 2,
				steps: []testStep{
					{p1, actionCard{}, offSeason},
					{p2, actionCard{}, offSeason},
					{p2, actionCard{}, campaignSeason},
					// p2 plays another card before p1 hits campaign readiness.
					{p2, actionCard{}, campaignSeason},
					{p1, actionCard{}, votingSeason},
				},
			},
			{
				name:         "voting season waits until all players vote",
				offSeasonLen: 3,
				steps: []testStep{
					{p1, actionCard{}, offSeason},
					{p1, actionCard{}, offSeason},
					{p1, actionCard{}, campaignSeason},
					{p2, actionCard{}, campaignSeason},
					{p2, actionCard{}, campaignSeason},
					{p2, actionCard{}, votingSeason},
					{p1, votingCard(""), votingSeason},
					{p2, votingCard(""), offSeason},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e := newElectionStateMachine(tt.offSeasonLen, []PID{p1, p2})
				for i, step := range tt.steps {
					got := e.HandleCardPlayed(step.pid, step.card)
					if got != step.wantSeason {
						t.Fatalf("step %d: HandleCardPlayed(%+v) got season %s but wanted %s", i, step, got, step.wantSeason)
					}
				}
			})
		}
	})
}

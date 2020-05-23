package couchcampaign

// import (
// 	"testing"
// )

// func TestElectionTracker(t *testing.T) {
// 	p1 := NewPID()
// 	p2 := NewPID()

// 	t.Run("2p", func(t *testing.T) {
// 		type testStep struct {
// 			pid        PID
// 			card       Card
// 			wantSeason electionSeason
// 		}
// 		tests := []struct {
// 			name         string
// 			offSeasonLen int
// 			steps        []testStep
// 		}{
// 			{
// 				name:         "not enough cards to trigger campaign season",
// 				offSeasonLen: 2,
// 				steps: []testStep{
// 					{p1, Card{}, offSeason},
// 					{p2, Card{}, offSeason},
// 				},
// 			},
// 			{
// 				name:         "campaign season triggered",
// 				offSeasonLen: 2,
// 				steps: []testStep{
// 					{p1, Card{}, offSeason},
// 					{p2, Card{}, offSeason},
// 					{p1, Card{}, campaignSeason},
// 					{p1, Card{}, campaignSeason},
// 				},
// 			},
// 			{
// 				name:         "voting season triggered",
// 				offSeasonLen: 2,
// 				steps: []testStep{
// 					{p1, Card{}, offSeason},
// 					{p2, Card{}, offSeason},
// 					{p2, Card{}, campaignSeason},
// 					// p2 plays another card before p1 hits campaign readiness.
// 					{p2, Card{}, campaignSeason},
// 					{p1, Card{}, votingSeason},
// 				},
// 			},
// 			{
// 				name:         "voting season waits until all players vote",
// 				offSeasonLen: 3,
// 				steps: []testStep{
// 					{p1, Card{}, offSeason},
// 					{p1, Card{}, offSeason},
// 					{p1, Card{}, campaignSeason},
// 					{p2, Card{}, campaignSeason},
// 					{p2, Card{}, campaignSeason},
// 					{p2, Card{}, votingSeason},
// 					{p1, TheVotingCard, votingSeason},
// 					{p2, TheVotingCard, offSeason},
// 				},
// 			},
// 		}
// 		for _, tt := range tests {
// 			t.Run(tt.name, func(t *testing.T) {
// 				e := newElectionStateMachine(tt.offSeasonLen, []PID{p1, p2})
// 				for i, step := range tt.steps {
// 					got := e.HandleCardPlayed(step.pid, step.card, nil)
// 					if got != step.wantSeason {
// 						t.Fatalf("step %d: HandleCardPlayed(%+v) got season %s but wanted %s", i, step, got, step.wantSeason)
// 					}
// 				}
// 			})
// 		}
// 	})
// }

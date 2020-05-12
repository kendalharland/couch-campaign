package couchcampaign

import (
	"github.com/gobuffalo/uuid"
)

type electionSeason int

const (
	offSeason electionSeason = iota
	campaignSeason
	votingSeason
)

func (s electionSeason) String() string {
	switch s {
	case offSeason:
		return "OffSeason"
	case campaignSeason:
		return "CampaignSeason"
	case votingSeason:
		return "VotingSeason"
	}
	return ""
}

type electionTracker struct {
	pids          []uuid.UUID
	numCardsUntil int

	currentSeason electionSeason
	numCardsSince map[uuid.UUID]int
	isVoting      map[uuid.UUID]bool
}

func newElectionTracker(offSeasonLength int, pids []uuid.UUID) *electionTracker {
	e := &electionTracker{
		pids:          pids,
		numCardsUntil: offSeasonLength,
	}
	e.reset()
	return e
}

func (e *electionTracker) CurrentSeason() electionSeason {
	return e.currentSeason
}

func (e *electionTracker) Update(input message) electionSeason {
	pid := input.pid
	switch e.currentSeason {
	case offSeason:
		// The off season ends as soon as at least one player activates
		// numCardsUntil action cards, which signifies the start of campaign
		// season.
		if _, ok := input.card.(actionCard); !ok {
			break
		}
		e.numCardsSince[pid]++
		if e.numCardsSince[pid] >= e.numCardsUntil {
			e.currentSeason = campaignSeason
		}
	case campaignSeason:
		// Campaign season ends as soon as the last player to activate
		// numCardsUntil action cards does so, which signifies the start of
		// voting season.
		if _, ok := input.card.(actionCard); !ok {
			break
		}
		e.numCardsSince[pid]++
		for _, count := range e.numCardsSince {
			if count < e.numCardsUntil {
				return e.currentSeason
			}
		}
		e.currentSeason = votingSeason
	case votingSeason:
		// Some players may still be looking at the card they drew right before
		// voting season, so the season ends as soon as all players are waiting
		// for votes to be counted.
		if _, ok := input.card.(votingCard); !ok {
			break
		}
		e.isVoting[pid] = true
		for _, voting := range e.isVoting {
			if !voting {
				return e.currentSeason
			}
		}
		// Everyone is voting, compute the results, reset every player's
		// progress and restart the election cycle.

		// TODO: Compute results.
		e.reset()
	}

	return e.currentSeason
}

func (e *electionTracker) reset() {
	e.currentSeason = offSeason
	e.numCardsSince = make(map[uuid.UUID]int)
	e.isVoting = make(map[uuid.UUID]bool)
	for _, pid := range e.pids {
		e.numCardsSince[pid] = 0
		e.isVoting[pid] = false
	}
}

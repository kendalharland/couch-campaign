package couchcampaign

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

type electionStateMachine struct {
	currentSeason electionSeason
	pids          []PID
	numCardsUntil int
	numCardsSince map[PID]int
	isVoting      map[PID]bool
}

func newElectionStateMachine(offSeasonLength int, pids []PID) *electionStateMachine {
	e := &electionStateMachine{
		pids:          pids,
		numCardsUntil: offSeasonLength,
	}
	e.reset()
	return e
}

func (e *electionStateMachine) CurrentSeason() electionSeason {
	return e.currentSeason
}

func (e *electionStateMachine) HandleCardPlayed(pid PID, card Card) electionSeason {
	t := cardTypeOf(card)

	switch e.currentSeason {
	case offSeason:
		// The off season ends as soon as at least one player activates
		// numCardsUntil action cards, which signifies the start of campaign
		// season.
		if t != actionCardType {
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
		if t != actionCardType {
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
		if t != votingCardType {
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

func (e *electionStateMachine) reset() {
	e.currentSeason = offSeason
	e.numCardsSince = make(map[PID]int)
	e.isVoting = make(map[PID]bool)
	for _, pid := range e.pids {
		e.numCardsSince[pid] = 0
		e.isVoting[pid] = false
	}
}

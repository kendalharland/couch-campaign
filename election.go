package couchcampaign

// type electionSeason int

// const (
// 	offSeason electionSeason = iota
// 	campaignSeason
// 	votingSeason
// )

// func (s electionSeason) String() string {
// 	switch s {
// 	case offSeason:
// 		return "OffSeason"
// 	case campaignSeason:
// 		return "CampaignSeason"
// 	case votingSeason:
// 		return "VotingSeason"
// 	}
// 	return ""
// }

// type electionStateMachine struct {
// 	currentSeason electionSeason
// 	pids          []PID
// 	numCardsUntil int
// 	numCardsSince map[PID]int
// 	isVoting      map[PID]bool
// }

// func newElectionStateMachine(offSeasonLength int, pids []PID) *electionStateMachine {
// 	e := &electionStateMachine{
// 		pids:          pids,
// 		numCardsUntil: offSeasonLength,
// 	}
// 	e.reset()
// 	return e
// }

// func (e *electionStateMachine) CurrentSeason() electionSeason {
// 	return e.currentSeason
// }

// func (e *electionStateMachine) HandleCardPlayed(pid PID, card Card, cs *CardService) electionSeason {
// 	t := cs.CardTypeOf(card)

// 	switch e.currentSeason {
// 	case offSeason:
// 		// The off season ends as soon as at least one player activates
// 		// numCardsUntil action cards, which signifies the start of campaign
// 		// season.
// 		if t != actionCardType {
// 			break
// 		}
// 		e.numCardsSince[pid]++
// 		if e.numCardsSince[pid] >= e.numCardsUntil {
// 			e.currentSeason = campaignSeason
// 		}
// 	case campaignSeason:
// 		// Campaign season ends as soon as the last player to activate
// 		// numCardsUntil action cards does so, which signifies the start of
// 		// voting season.
// 		if t != actionCardType {
// 			break
// 		}
// 		e.numCardsSince[pid]++
// 		for _, count := range e.numCardsSince {
// 			if count < e.numCardsUntil {
// 				return e.currentSeason
// 			}
// 		}
// 		e.currentSeason = votingSeason
// 	}
// 	return e.currentSeason
// }

// func (e *electionStateMachine) HandleIsVoting(pid PID) {
// 	for _, opponent := range e.pids {
// 		if opponent == pid {
// 			continue
// 		}
// 		for _, voting := range e.isVoting {
// 			if !voting {
// 				return
// 			}
// 		}
// 	}
// 	e.reset()
// }

// func (e *electionStateMachine) reset() {
// 	e.currentSeason = offSeason
// 	e.numCardsSince = make(map[PID]int)
// 	e.isVoting = make(map[PID]bool)
// 	for _, pid := range e.pids {
// 		e.numCardsSince[pid] = 0
// 		e.isVoting[pid] = false
// 	}
// }

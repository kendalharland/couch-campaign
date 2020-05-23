package couchcampaign

import (
	"couchcampaign/cardgame"
)

type PID = cardgame.PID
type ClientWorker = cardgame.ClientWorker
type CardService = cardgame.CardService

func NewGame(clients map[PID]*ClientWorker) (*cardgame.Game, error) {
	pids := make([]PID, 0, len(clients))
	for pid := range clients {
		pids = append(pids, pid)
	}

	state := cardgame.NewGameState(pids, &CardService{Cards: allCards})

	g := &cardgame.Game{
		state:          state,
		clients:        clients,
		clientJobs:     make(map[PID]chan cardgame.ClientJob),
		clientMessages: make(chan cardgame.ClientMessage, len(clients)),
		clientErrors:   make(chan cardgame.ClientError, len(clients)),
	}

	for _, pid := range g.state.pids {
		g.clientJobs[pid] = make(chan cardgame.ClientJob, 2)
	}

	return g, nil
}

package couchcampaign

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Game drives a server-side Game session.
//
// Its main purpose is to route messages to and from clients and monitor each
// client's connection state.
type Game struct {
	state          *GameState
	clients        map[PID]*ClientWorker
	clientJobs     map[PID]chan ClientJob
	clientMessages chan ClientMessage
	clientErrors   chan ClientError
	mu             sync.Mutex
}

func NewGame(clients map[PID]*ClientWorker) (*Game, error) {
	baseCards, err := loadBaseCards()
	if err != nil {
		return nil, fmt.Errorf("loadBaseCards: %w", err)
	}

	pids := make([]PID, 0, len(clients))
	for pid := range clients {
		pids = append(pids, pid)
	}

	state := newGameState(pids, baseCards)

	g := &Game{
		state:          state,
		clients:        clients,
		clientJobs:     make(map[PID]chan ClientJob),
		clientMessages: make(chan ClientMessage, len(clients)),
		clientErrors:   make(chan ClientError, len(clients)),
	}

	for _, pid := range g.state.pids {
		g.clientJobs[pid] = make(chan ClientJob, 2)
	}

	return g, nil
}

func (g *Game) Run(_ context.Context) {
	defer g.shutdown()

	// Spawn all clients.
	for _, pid := range g.state.pids {
		go g.clients[pid].Run(g.clientJobs[pid], g.clientMessages, g.clientErrors)
		g.clientJobs[pid] <- g.state.getNextJob(pid)
	}

	for !g.isOver() {
		select {
		case message := <-g.clientMessages:
			pids, err := g.handleMessage(message)
			if err != nil {
				log.Fatal(err)
			}
			for _, pid := range pids {
				g.clientJobs[pid] <- g.state.getNextJob(pid)
			}
		case err := <-g.clientErrors:
			if IsConnectionCloseError(err) {
				g.disconnect(err.PID)
				log.Printf("disconnected from client %v", err.PID)
				break
			}
			log.Printf("%v: %v", err.PID, err.Error())
		}
	}

	log.Println("Game over")
}

func (g *Game) handleMessage(m ClientMessage) (pids []PID, err error) {
	// TODO: Put this behind a flag.
	defer g.state.debugDumpDecks()

	switch m.Input {
	case DismissInfoCardInput:
		g.state.OnPlayerDismissedInfoCard(m.PID)
		return []PID{m.PID}, nil
	case AcceptActionCardInput:
		return g.state.OnPlayerAcceptedActionCard(m.PID)
	case RejectActionCardInput:
		return g.state.OnPlayerRejectedActionCard(m.PID)
	default:
		return g.state.updateElectionState(m.PID)
	}
}

func (g *Game) shutdown() {
	close(g.clientMessages)
	close(g.clientErrors)
	for _, pid := range g.state.pids {
		g.disconnect(pid)
	}
}

func (g *Game) disconnect(pid PID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	close(g.clientJobs[pid])
	delete(g.clientJobs, pid)

	delete(g.clients, pid)
	g.state.removePlayer(pid)
	for i, id := range g.state.pids {
		if id == pid {
			g.state.pids = append(g.state.pids[i:], g.state.pids[:i+1]...)
		}
	}
}

func (g *Game) isOver() bool {
	return len(g.clients) <= 0
}

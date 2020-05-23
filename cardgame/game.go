package cardgame

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
	logic          GameLogic
	clients        map[PID]*ClientWorker
	clientJobs     map[PID]chan ClientJob
	clientMessages chan ClientMessage
	clientErrors   chan ClientError
	mu             sync.Mutex
}

func (g *Game) Run(_ context.Context) {
	defer g.shutdown()

	// Spawn all clients.
	for _, pid := range g.logic.pids {
		go g.clients[pid].Run(g.clientJobs[pid], g.clientMessages, g.clientErrors)
		g.clientJobs[pid] <- g.logic.getNextJob(pid)
	}

	defer log.Println("Game over")

	for {
		select {
		case message := <-g.clientMessages:
			if err := g.handleMessage(message); err != nil {
				log.Fatal(err)
			}
		case err := <-g.clientErrors:
			if IsConnectionCloseError(err) {
				g.disconnect(err.PID)
				log.Printf("disconnected from client %v", err.PID)
				break
			}
			if len(g.clients) == 0 {
				return
			}
			log.Printf("%v: %v", err.PID, err.Error())
		}
	}
}

func (g *Game) handleMessage(m ClientMessage) error {
	var pids []PID
	var err error

	switch m.Input {
	case NoInput:
		pids, err = g.logic.OnCardShown(m.PID)
	case AcceptCardInput:
		pids, err = g.logic.OnCardAccepted(m.PID)
	case RejectCardInput:
		pids, err = g.logic.OnCardRejected(m.PID)
	case FailedInput:
		err = fmt.Errorf("failed to get input from player %v", m.PID)
	default:
		err = fmt.Errorf("invalid card input: %v", m.Input)
	}
	if err != nil {
		return err
	}

	for _, pid := range pids {
		g.clientJobs[pid] <- g.logic.getNextJob(pid)
	}

	return nil
}

func (g *Game) shutdown() {
	close(g.clientMessages)
	close(g.clientErrors)
	for _, pid := range g.logic.pids {
		g.disconnect(pid)
	}
}

func (g *Game) disconnect(pid PID) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.logic.removePlayer(pid)
	for i, id := range g.logic.pids {
		if id == pid {
			g.logic.pids = append(g.logic.pids[i:], g.state.pids[:i+1]...)
		}
	}

	close(g.clientJobs[pid])
	delete(g.clientJobs, pid)
	delete(g.clients, pid)
}

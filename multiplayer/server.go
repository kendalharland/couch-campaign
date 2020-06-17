package multiplayer

import (
	"log"
	"sync"
)

// TODO: Handle disconnects.
// TODO: Handle websocket message type.

// Server drives a server-side game session.
//
// It spawns as many "client" go-routines as necessary to handle all client connections.
// Each client go-routine writes messages recieved from the remote peer into the server's
// incoming message queue. The current go-routine processes all messages from the incoming
// queue by sending them to the Game implementation.
//
// Each remote peer is identified by a unique CID. The Game implementation can tell the
// server to send a message to a peer by returning a Message containing that peer's CID
// from HandleMessage.
type Server struct {
	clients          map[CID]*Client
	outgoingMessages map[CID]chan Message
	incomingMessages chan Message
	clientErrors     chan ClientError
	mu               sync.Mutex
}

// NewServer creates a new Server with the given clients.
func NewServer(clients map[CID]*Client) *Server {
	pids := make([]CID, 0, len(clients))
	for pid := range clients {
		pids = append(pids, pid)
	}

	s := &Server{
		clients:          clients,
		outgoingMessages: make(map[CID]chan Message),
		incomingMessages: make(chan Message, len(clients)),
		clientErrors:     make(chan ClientError, len(clients)),
	}

	for _, pid := range pids {
		s.outgoingMessages[pid] = make(chan Message, 2)
	}
	return s
}

// Run starts the game loop.
func (s *Server) Run(g Game) {
	defer s.shutdown()

	for cid, client := range s.clients {
		go client.Run(s.outgoingMessages[cid], s.incomingMessages, s.clientErrors)
	}

	for {
		select {
		case message := <-s.incomingMessages:
			responses, err := g.HandleMessage(message)
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range responses {
				s.outgoingMessages[r.CID] <- r
			}
		case err := <-s.clientErrors:
			if err := g.HandleError(err); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// disconnect disconnects the given client from the game.
func (s *Server) disconnect(pid CID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.outgoingMessages[pid])
	delete(s.outgoingMessages, pid)
	delete(s.clients, pid)
}

// Shutdown disposes of this game context.
func (s *Server) shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.incomingMessages)
	close(s.clientErrors)
}

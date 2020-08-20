package multiplayer

import (
	"log"
	"sync"
)

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
	clients map[CID]*Client
	outs    map[CID]chan Message
	ins     chan Message
	errs    chan ClientError
	mu      sync.Mutex
}

// NewServer creates a new Server with the given clients.
func NewServer() *Server {
	return &Server{
		clients: make(map[CID]*Client),
		outs:    make(map[CID]chan Message),
		ins:     make(chan Message),
		errs:    make(chan ClientError),
	}
}

func (s *Server) AddClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: Handle disconnect

	s.clients[client.ID()] = client
	s.outs[client.ID()] = make(chan Message, 2)
}

func (s *Server) NClients() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.clients)
}

func (s *Server) Send(m Message) {
	s.outs[m.CID] <- m
}

// Run starts the game loop.
func (s *Server) Run(build GameBuilder) error {
	defer s.shutdown()

	ids := make([]CID, 0, len(s.clients))
	for id, client := range s.clients {
		ids = append(ids, id)
		go client.Run(s.outs[id], s.ins, s.errs)
	}

	game, err := build(ids)
	if err != nil {
		return err
	}

	for {
		select {
		case message := <-s.ins:
			if err := game.HandleInput(message.CID, message.Data); err != nil {
				log.Fatal(err)
			}
			// TODO: Implement server sent events.
		}
	}

	// TODO: Disconnect all remaining players and shutdown when game is over.
}

func (s *Server) disconnect(id CID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.outs[id])
	delete(s.outs, id)
	delete(s.clients, id)
}

func (s *Server) shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

}

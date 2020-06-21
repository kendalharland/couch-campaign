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
func NewServer() *Server {
	return &Server{
		clients:          make(map[CID]*Client),
		outgoingMessages: make(map[CID]chan Message),
		incomingMessages: make(chan Message),
		clientErrors:     make(chan ClientError),
	}
}

func (s *Server) AddClient(cid CID, client *Client) {
	s.clients[cid] = client
	s.outgoingMessages[cid] = make(chan Message, 2)
}

// Run starts the game loop.
func (s *Server) Run(g Game) error {
	defer s.shutdown()

	for cid, client := range s.clients {
		if err := g.AddPlayer(cid); err != nil {
			return err
		}
		go client.Run(s.outgoingMessages[cid], s.incomingMessages, s.clientErrors)
	}

	if err := g.Start(); err != nil {
		return err
	}

	for {
		select {
		case message := <-s.incomingMessages:
			if err := g.HandleInput(message.CID, message.Data); err != nil {
				log.Fatal(err)
			}
		case message := <-g.Outputs():
			s.outgoingMessages[message.CID] <- message
		case err := <-s.clientErrors:
			if err := g.HandleError(err); err != nil {
				return err
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

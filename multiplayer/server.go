package multiplayer

import (
	"fmt"
	"log"
	"sync"

	"github.com/gobuffalo/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	pb "couchcampaign/proto"
)

var (
	lobbyStateMessage,
	runningStateMessage []byte
)

func init() {
	prot := &pb.Message{Content: &pb.Message_SessionState{SessionState: pb.SessionState_LOBBY}}
	data, err := proto.Marshal(prot)
	if err != nil {
		log.Fatalf("proto.Marshal(%+v): %v", prot, err)
	}
	lobbyStateMessage = data

	prot = &pb.Message{Content: &pb.Message_SessionState{SessionState: pb.SessionState_RUNNING}}
	data, err = proto.Marshal(prot)
	if err != nil {
		log.Fatalf("proto.Marshal(%+v): %v", prot, err)
	}
	runningStateMessage = data
}

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

func (s *Server) AddConnection(c *websocket.Conn) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}
	id := CID(uuid.String())
	s.AddClient(NewClient(id, c))
	return nil
}

func (s *Server) AddClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.ID()] = client
	s.outs[client.ID()] = make(chan Message, 2)
	s.outs[client.ID()] <- Message{
		CID:          client.ID(),
		Data:         lobbyStateMessage,
		SkipResponse: true,
	}
	// TODO: Handle disconnects.
}

func (s *Server) NClients() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.clients)
}

// Run starts the game loop.
func (s *Server) Run(build GameBuilder) error {
	defer s.shutdown()

	ids := make([]CID, 0, len(s.clients))
	for id, client := range s.clients {
		ids = append(ids, id)
		s.outs[id] <- Message{
			CID:          id,
			Data:         runningStateMessage,
			SkipResponse: true,
		}
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

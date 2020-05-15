package main

import (
	"couchcampaign"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"couchcampaign/multiplayer"
	pb "couchcampaign/proto"
)

var upgrader = websocket.Upgrader{} // use default options

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type GameServer struct {
	game  *multiplayer.Server
	conns map[multiplayer.CID]*websocket.Conn

	mu sync.Mutex
}

func NewGameServer() *GameServer {
	return &GameServer{
		conns: make(map[multiplayer.CID]*websocket.Conn),
	}
}

func (s *GameServer) InstallHandlers(r *mux.Router) {
	r.HandleFunc("/", s.status)
	r.HandleFunc("/connect", s.connect)
	r.HandleFunc("/start", s.start)
	r.HandleFunc("/socket", s.socket)
}

func (s *GameServer) ensureNotStarted() error {
	if s.game != nil {
		return errors.New("game has already started")
	}
	return nil
}

func (s *GameServer) status(w http.ResponseWriter, r *http.Request) {
	couchcampaign.Respond(w, http.StatusOK, "game is running")
}

func (s *GameServer) start(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer r.Body.Close()

	if err := s.ensureNotStarted(); err != nil {
		couchcampaign.RespondWithError(w, err)
	}

	clients := make(map[multiplayer.CID]*multiplayer.Client)
	cids := make([]multiplayer.CID, 0, len(s.conns))
	for cid, conn := range s.conns {
		clients[cid] = multiplayer.NewClient(cid, conn)
		cids = append(cids, cid)
	}

	server := multiplayer.NewServer(clients)
	game := couchcampaign.NewGameWithCIDs(cids)
	go server.Run(game)

	log.Println("game started")
	couchcampaign.Respond(w, http.StatusOK, "1")

	message := pb.Message{Content: &pb.Message_SessionState{SessionState: pb.SessionState_RUNNING}}
	payload, err := proto.Marshal(&message)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
	}
	for _, conn := range s.conns {
		if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
			log.Println(err)
		}
	}
}

func (s *GameServer) connect(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer r.Body.Close()

	if err := s.ensureNotStarted(); err != nil {
		couchcampaign.RespondWithError(w, err)
	}

	id := multiplayer.CID("TODO")

	// Upgrade to a websocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	s.conns[id] = conn

	message := pb.Message{Content: &pb.Message_SessionState{SessionState: pb.SessionState_LOBBY}}
	payload, err := proto.Marshal(&message)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		log.Println(err)
	}
}

func (s *GameServer) socket(w http.ResponseWriter, r *http.Request) {
	socketAddr := url.URL{
		Scheme: "ws",
		Host:   r.Host,
		Path:   "connect",
	}
	couchcampaign.Respond(w, http.StatusOK, socketAddr.String())
}

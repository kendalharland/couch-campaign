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

var (
	errGameStarted      = errors.New("game is already started")
	errNotEnoughPlayers = errors.New("game does not have enough players")
	errMaxPlayers       = errors.New("game is full")
)

const (
	minPlayerCount = 0
	maxPlayerCount = 50
)

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type GameServer struct {
	isGameRunning bool
	conns         map[multiplayer.CID]*websocket.Conn

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

func (s *GameServer) start(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer r.Body.Close()

	if s.isGameRunning {
		couchcampaign.RespondWithError(w, errGameStarted)
		return
	}
	if len(s.conns) < minPlayerCount {
		couchcampaign.RespondWithError(w, errNotEnoughPlayers)
		return
	}

	server := multiplayer.NewServer()
	for cid, conn := range s.conns {
		server.AddClient(cid, multiplayer.NewClient(cid, conn))
	}

	go server.Run(couchcampaign.NewGame())
	couchcampaign.Respond(w, http.StatusOK, "")
}

func (s *GameServer) connect(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer r.Body.Close()

	if s.isGameRunning {
		couchcampaign.RespondWithError(w, errGameStarted)
		return
	}
	if len(s.conns) >= maxPlayerCount {
		couchcampaign.RespondWithError(w, errMaxPlayers)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}

	id := multiplayer.CID("TODO_let_the_client_set_this")
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

func (s *GameServer) status(w http.ResponseWriter, r *http.Request) {
	couchcampaign.Respond(w, http.StatusOK, "game is running")
}

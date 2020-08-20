package main

import (
	"couchcampaign"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"couchcampaign/multiplayer"
)

var (
	upgrader     = websocket.Upgrader{} // use default options
	nextClientID = 0
)

var (
	errGameStarted      = errors.New("game is already started")
	errNotEnoughPlayers = errors.New("game does not have enough players")
	errMaxPlayers       = errors.New("game is full")
)

const (
	minPlayerCount = 1
	maxPlayerCount = 50
)

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type GameServer struct {
	isGameRunning bool
	server        *multiplayer.Server

	mu sync.Mutex
}

func NewGameServer() *GameServer {
	return &GameServer{
		server: multiplayer.NewServer(),
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
	if s.server.NClients() < minPlayerCount {
		couchcampaign.RespondWithError(w, errNotEnoughPlayers)
		return
	}

	go s.server.Run(multiplayer.GameBuilder(couchcampaign.NewGame))

	// TODO: Broadcast session started.
	couchcampaign.Respond(w, http.StatusOK, "")
	s.isGameRunning = true
}

func (s *GameServer) connect(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer r.Body.Close()

	if s.isGameRunning {
		couchcampaign.RespondWithError(w, errGameStarted)
		return
	}
	if s.server.NClients() >= maxPlayerCount {
		couchcampaign.RespondWithError(w, errMaxPlayers)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	if err := s.server.AddConnection(conn); err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	couchcampaign.Respond(w, http.StatusOK, "")
}

func (s *GameServer) socket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	socketAddr := url.URL{
		Scheme: "ws",
		Host:   r.Host,
		Path:   "connect",
	}
	couchcampaign.Respond(w, http.StatusOK, socketAddr.String())
}

func (s *GameServer) status(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	couchcampaign.Respond(w, http.StatusOK, fmt.Sprintf(`Players: (%d/%d)`, s.server.NClients(), maxPlayerCount))
}

// 	for id, state := range g.Ctx.SnapshotPlayerStates() {
// 		playerStateMessage, err := proto.Marshal(playerStateToMessageProto(state))
// 		if err != nil {
// 			return err
// 		}
// 		// Alert the client that the session is now running.
// 		g.outputs <- multiplayer.Message{CID: id, Data: sessionStartedMessage, SkipResponse: true}

func SendSessionMessage()

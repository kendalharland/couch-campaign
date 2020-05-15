package main

import (
	"context"
	"couchcampaign"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/golang/protobuf/jsonpb"

	"github.com/gorilla/mux"

	pb "couchcampaign/proto"
)

// LobbyServer handles HTTP requests for game games.
type LobbyServer struct {
	version string
	impl    *LobbyServerImpl
}

// NewLobbyServer creates a new LobbyServer.
func NewLobbyServer(version string, minPort, maxPort int) *LobbyServer {
	return &LobbyServer{
		version: version,
		impl: &LobbyServerImpl{
			minPort: minPort,
			maxPort: maxPort,
		},
	}
}

// InstallHandlers registers lobby HTTP handlers to mux.
func (s *LobbyServer) InstallHandlers(r *mux.Router) {
	r.HandleFunc("/", s.index)
	r.HandleFunc("/lobby/create", s.createGame)
	r.HandleFunc("/lobby/list", s.listGames)
	r.HandleFunc("/lobby/join", s.joinGame)
}

func (s *LobbyServer) index(w http.ResponseWriter, r *http.Request) {
	couchcampaign.Respond(w, http.StatusOK, fmt.Sprintf("couchcampaign version %s", s.version))
}

func (s *LobbyServer) createGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := &pb.CreateGameRequest{}
	if err := jsonpb.Unmarshal(r.Body, req); err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	res, err := s.impl.CreateGame(context.Background(), req)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
	}
	couchcampaign.Respond(w, http.StatusOK, string(data))
}

func (s *LobbyServer) joinGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := &pb.JoinGameRequest{}
	if err := jsonpb.Unmarshal(r.Body, req); err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	res, err := s.impl.JoinGame(context.Background(), req)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
	}
	couchcampaign.Respond(w, http.StatusOK, string(data))
}

func (s *LobbyServer) listGames(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	req := &pb.ListGamesRequest{}
	if err := jsonpb.Unmarshal(r.Body, req); err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	res, err := s.impl.ListGames(context.Background(), req)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		couchcampaign.RespondWithError(w, err)
	}
	couchcampaign.Respond(w, http.StatusOK, string(data))
}

type LobbyServerImpl struct {
	games            []*Game
	minPort, maxPort int
}

// Game is a staging area for a game that has not yet started.
type Game struct {
	ID   string `json:"id"`
	Port int
}

// CreateGame creates a new lobby on this server.
func (s *LobbyServerImpl) CreateGame(ctx context.Context, m *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	port := s.minPort
	for _, lobby := range s.games {
		log.Println(lobby.Port)
		if lobby.Port == port {
			port++
		}
	}
	if port > s.maxPort {
		return nil, errors.New("no ports available")
	}

	id := fmt.Sprintf("%d", rand.Intn(999999))
	s.games = append(s.games, &Game{
		ID:   id,
		Port: port,
	})

	// Set the Cmd explicitly because exec.Command() cannot find binaries
	// with a suffix on windows.
	cmd := exec.Cmd{
		Path:   "couchcampaign",
		Args:   []string{"couchcampaign", "-port", fmt.Sprintf("%d", port)},
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	log.Println(cmd.String())
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	gameURL := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
	}
	return &pb.CreateGameResponse{GameUrl: gameURL.String()}, nil
}

// ListGames returns the games on this server.
func (s *LobbyServerImpl) ListGames(ctx context.Context, m *pb.ListGamesRequest) (*pb.ListGamesResponse, error) {
	res := &pb.ListGamesResponse{}
	for _, game := range s.games {
		res.Games = append(res.Games, &pb.GameInfo{
			Id: game.ID,
		})
	}
	return res, nil
}

// JoinGame adds a player to a lobby.
func (s *LobbyServerImpl) JoinGame(ctx context.Context, m *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	if m.GameId == "" {
		return nil, couchcampaign.Errorf(http.StatusBadRequest, "A lobby ID is required")
	}

	lpos := -1
	for i, lobby := range s.games {
		if lobby.ID == m.GameId {
			lpos = i
			break
		}
	}
	if lpos < 0 {
		return nil, couchcampaign.Errorf(http.StatusNotFound, "Game %q not found", m.GameId)
	}
	l := s.games[lpos]
	return &pb.JoinGameResponse{
		GameUrl: fmt.Sprintf("ws://localhost:%d", l.Port),
	}, nil
}

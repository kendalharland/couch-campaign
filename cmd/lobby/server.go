package main

import (
	"context"
	"couchcampaign"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/golang/protobuf/jsonpb"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// TODO: Can LobbyServer be functions instead of struct?
// TODO: Include healthcheck URL for games so that LobbyServer can clean up if the process dies.

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

	req := &couchcampaign.CreateGameRequest{}
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

	req := &couchcampaign.JoinGameRequest{}
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

	req := &couchcampaign.ListGamesRequest{}
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
func (s *LobbyServerImpl) CreateGame(ctx context.Context, m *couchcampaign.CreateGameRequest) (*couchcampaign.CreateGameResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	port := s.minPort
	for _, lobby := range s.games {
		if lobby.Port == port {
			port++
		}
	}
	if port > s.maxPort {
		return nil, errors.New("no ports available")
	}

	s.games = append(s.games, &Game{
		ID:   id.String(),
		Port: port,
	})

	cmd := exec.Command("couchcampaign", fmt.Sprintf("-port=%d", port))
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// TODO: Find a way to pass the original URL to this method instead of
	// hardcoding 'localhost'.
	gameURL := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", port),
	}
	return &couchcampaign.CreateGameResponse{GameUrl: gameURL.String()}, nil
}

// ListGames returns the games on this server.
func (s *LobbyServerImpl) ListGames(ctx context.Context, m *couchcampaign.ListGamesRequest) (*couchcampaign.ListGamesResponse, error) {
	res := &couchcampaign.ListGamesResponse{}
	for _, game := range s.games {
		res.Games = append(res.Games, &couchcampaign.GameInfo{
			Id: game.ID,
		})
	}
	return res, nil
}

// JoinGame adds a player to a lobby.
func (s *LobbyServerImpl) JoinGame(ctx context.Context, m *couchcampaign.JoinGameRequest) (*couchcampaign.JoinGameResponse, error) {
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
	return &couchcampaign.JoinGameResponse{
		GameUrl: fmt.Sprintf("ws://localhost:%d", l.Port),
	}, nil
}

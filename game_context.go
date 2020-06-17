package couchcampaign

import (
	"fmt"
)

type gameContext struct {
	players map[string]*PlayerState
}

func (g *gameContext) PlayerIDs() []string {
	ids := make([]string, 0, len(g.players))
	for id := range g.players {
		ids = append(ids, id)
	}
	return ids
}

func (g *gameContext) SnapshotPlayerStates() map[string]PlayerState {
	snapshot := make(map[string]PlayerState)
	for k, v := range g.players {
		snapshot[k] = *v
	}
	return snapshot
}

func (g *gameContext) AddPlayer(id string) error {
	if _, ok := g.players[id]; ok {
		return fmt.Errorf("player with id %q already exists", id)
	}
	g.players[id] = &PlayerState{ID: id}
	return nil
}

func (g *gameContext) GetPlayer(id string) (*PlayerState, error) {
	ps, ok := g.players[id]
	if !ok {
		return nil, fmt.Errorf("no player with id %q", id)
	}
	clone := *ps
	return &clone, nil
}

func (g *gameContext) SetPlayer(id string, state PlayerState) error {
	if _, ok := g.players[id]; !ok {
		return fmt.Errorf("no player with id %q", id)
	}
	g.players[id] = &state
	return nil
}

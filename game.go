package couchcampaign

import (
	"errors"
	"fmt"
	"log"

	"couchcampaign/multiplayer"
	couchcampaignpb "couchcampaign/proto"

	"google.golang.org/protobuf/proto"
)

const devScriptsPath = "scripts/main.star"

func NewGame() *game {
	return &game{
		ctx: &gameContext{
			players: make(map[string]*PlayerState),
		},
	}
}

type game struct {
	isStarted bool
	ctx       *gameContext
	scripts   *scriptManager
}

func (g *game) AddPlayer(cid multiplayer.CID) error {
	if g.isStarted {
		return errors.New("game is already started")
	}
	g.ctx.AddPlayer(string(cid))
	return nil
}

func (g *game) Start() error {
	if g.isStarted {
		return errors.New("game is already started")
	}
	g.scripts = newScriptManager(g.ctx)
	if err := g.scripts.LoadMainScript(devScriptsPath); err != nil {
		return err
	}
	g.isStarted = true
	return nil
}

func (g *game) HandleMessage(m multiplayer.Message) ([]multiplayer.Message, error) {
	if !g.isStarted {
		return nil, errors.New("game is not started")
	}

	input, err := parseInput(m.Data)
	if err != nil {
		return nil, fmt.Errorf("parseInput: %w", err)
	}

	oldPlayerStates := g.ctx.SnapshotPlayerStates()
	if err := g.scripts.HandleInput(m.CID, input); err != nil {
		return nil, fmt.Errorf("HandleInput: %w", err)
	}

	var messages []multiplayer.Message
	for id, state := range g.ctx.SnapshotPlayerStates() {
		if oldPlayerStates[id] == state {
			continue
		}
		data, err := proto.Marshal(playerStateToProto(state))
		if err != nil {
			return nil, err
		}
		messages = append(messages, multiplayer.Message{
			CID:  id,
			Data: data,
		})
	}
	return messages, nil
}

// HandleError handles a client error.
func (g *game) HandleError(e multiplayer.ClientError) error {
	if !g.isStarted {
		return errors.New("game is not started")
	}
	log.Printf("error: %v: %v", e.CID, e.Err)
	return nil
}

// Close disposes this game and all of its resources.
func (g *game) Close() {}

func parseInput(input []byte) (Input, error) {
	value := string(input)
	switch value {
	case "accept":
		return InputCardAccepted, nil
	case "reject":
		return InputCardRejected, nil
	case "show":
		return InputCardShown, nil
	default:
		return InputErr, fmt.Errorf("invalid input: %q", value)
	}
}

func playerStateToProto(ps PlayerState) *couchcampaignpb.PlayerState {
	return &couchcampaignpb.PlayerState{
		Id:                 ps.ID,
		Leader:             ps.Leader,
		LeaderTimeInOffice: int32(ps.LeaderTimeInOffice),
		Wealth:             int32(ps.Wealth),
		Health:             int32(ps.Health),
		Stability:          int32(ps.Stability),
		Card: &couchcampaignpb.Card{
			Text:       ps.CardText,
			Speaker:    ps.CardSpeaker,
			AcceptText: ps.CardAcceptText,
			RejectText: ps.CardRejectText,
		},
	}
}

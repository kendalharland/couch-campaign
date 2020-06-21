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
		outputs: make(chan multiplayer.Message, 2),
	}
}

type game struct {
	ctx       *gameContext
	outputs   chan multiplayer.Message
	isStarted bool
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

	sessionStartedMessage, err := proto.Marshal(&couchcampaignpb.Message{
		Content: &couchcampaignpb.Message_SessionState{
			SessionState: couchcampaignpb.SessionState_RUNNING,
		},
	})
	if err != nil {
		return err
	}

	for id, state := range g.ctx.SnapshotPlayerStates() {
		playerStateMessage, err := proto.Marshal(playerStateToMessageProto(state))
		if err != nil {
			return err
		}
		// Alert the client that the session is now running.
		g.outputs <- multiplayer.Message{CID: id, Data: sessionStartedMessage, SkipResponse: true}

		// Send the client's initial state.
		g.outputs <- multiplayer.Message{CID: id, Data: playerStateMessage}
	}

	g.isStarted = true
	return nil
}

func (g *game) Stop() error { return nil }

func (g *game) HandleInput(cid multiplayer.CID, data []byte) error {
	if !g.isStarted {
		return errors.New("game is not started")
	}

	input, err := parseInput(data)
	if err != nil {
		return fmt.Errorf("parseInput: %w", err)
	}

	oldPlayerStates := g.ctx.SnapshotPlayerStates()
	if err := g.scripts.HandleInput(cid, input); err != nil {
		return fmt.Errorf("HandleInput: %w", err)
	}

	for id, state := range g.ctx.SnapshotPlayerStates() {
		if oldPlayerStates[id] == state {
			continue
		}
		data, err := proto.Marshal(playerStateToMessageProto(state))
		if err != nil {
			return err
		}
		g.outputs <- multiplayer.Message{CID: id, Data: data}
	}
	return nil
}

func (g *game) Outputs() <-chan multiplayer.Message {
	return g.outputs
}

// HandleError handles a client error.
func (g *game) HandleError(e multiplayer.ClientError) error {
	if !g.isStarted {
		return errors.New("game is not started")
	}
	log.Printf("error: %v: %v", e.CID, e.Err)
	return nil
}

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

func playerStateToMessageProto(state PlayerState) *couchcampaignpb.Message {
	return &couchcampaignpb.Message{
		Content: &couchcampaignpb.Message_PlayerState{
			PlayerState: playerStateToProto(state),
		},
	}
}

func playerStateToProto(state PlayerState) *couchcampaignpb.PlayerState {
	return &couchcampaignpb.PlayerState{
		Id:                 state.ID,
		Leader:             state.Leader,
		LeaderTimeInOffice: int32(state.LeaderTimeInOffice),
		Wealth:             int32(state.Wealth),
		Health:             int32(state.Health),
		Stability:          int32(state.Stability),
		Card: &couchcampaignpb.Card{
			Text:       state.CardText,
			Speaker:    state.CardSpeaker,
			AcceptText: state.CardAcceptText,
			RejectText: state.CardRejectText,
		},
	}
}

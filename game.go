package couchcampaign

import (
	"fmt"

	"couchcampaign/multiplayer"
	couchcampaignpb "couchcampaign/proto"
	"couchcampaign/starlarkgame"

	"google.golang.org/protobuf/proto"
)

const scriptsFilename = "scripts/main.star"

type Game struct {
	g *starlarkgame.Game
}

func NewGame(playerIDs []multiplayer.CID) (*Game, error) {
	ctx := starlarkgame.NewContext()
	for _, id := range playerIDs {
		ctx.AddPlayer(id)
	}
	g, err := starlarkgame.New(ctx, scriptsFilename)
	if err != nil {
		return nil, err
	}
	return &Game{g}, nil
}

func (g *Game) HandleInput(cid multiplayer.CID, data []byte) error {
	input, err := parseInput(data)
	if err != nil {
		return fmt.Errorf("parseInput: %w", err)
	}
	if err := g.g.HandleInput(cid, string(input)); err != nil {
		return fmt.Errorf("HandleInput(%v, %v): %w", cid, input, err)
	}
	return nil
}

func (g *Game) GetPlayerState(cid multiplayer.CID) ([]byte, error) {
	state, err := g.g.GetPlayerState(cid)
	if err != nil {
		return nil, fmt.Errorf("GetPlayerState: %w", err)
	}
	data, err := proto.Marshal(playerStateToMessageProto(*state))
	if err != nil {
		return nil, fmt.Errorf("proto.Marshal: %w", err)
	}
	return data, nil
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

func playerStateToMessageProto(state starlarkgame.PlayerState) *couchcampaignpb.Message {
	return &couchcampaignpb.Message{
		Content: &couchcampaignpb.Message_PlayerState{
			PlayerState: playerStateToProto(state),
		},
	}
}

func playerStateToProto(state starlarkgame.PlayerState) *couchcampaignpb.PlayerState {
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

// oldPlayerStates := g.Ctx.SnapshotPlayerStates()
// 	if err := g.ScriptManager.HandleInput(cid, input); err != nil {
// 		return fmt.Errorf("HandleInput: %w", err)
// 	}

// 	for id, state := range g.Ctx.SnapshotPlayerStates() {
// 		if oldPlayerStates[id] == state {
// 			continue
// 		}
// 		data, err := proto.Marshal(playerStateToMessageProto(state))
// 		if err != nil {
// 			return err
// 		}
// 		g.outputs <- multiplayer.Message{CID: id, Data: data}
// 	}
// 	return nil

// func (g *Game) Start() error {
// 	g.ScriptManager = newScriptManager(g.Ctx)
// 	if err := g.ScriptManager.LoadMainScript(scriptsFilename); err != nil {
// 		return err
// 	}

// 	sessionStartedMessage, err := proto.Marshal(&couchcampaignpb.Message{
// 		Content: &couchcampaignpb.Message_SessionState{
// 			SessionState: couchcampaignpb.SessionState_RUNNING,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for id, state := range g.Ctx.SnapshotPlayerStates() {
// 		playerStateMessage, err := proto.Marshal(playerStateToMessageProto(state))
// 		if err != nil {
// 			return err
// 		}
// 		// Alert the client that the session is now running.
// 		g.outputs <- multiplayer.Message{CID: id, Data: sessionStartedMessage, SkipResponse: true}

// 		// Send the client's initial state.
// 		g.outputs <- multiplayer.Message{CID: id, Data: playerStateMessage}
// 	}

// 	return nil
// }

// func (g *Game) Outputs() <-chan multiplayer.Message {
// 	return g.outputs
// }

// // HandleError handles a client error.
// func (g *Game) HandleError(e multiplayer.ClientError) error {
// 	log.Printf("error: %v: %v", e.CID, e.Err)
// 	return nil
// }

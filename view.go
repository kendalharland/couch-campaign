package couchcampaign

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

func renderCard(c Card, s stats) ([]byte, error) {
	m := PlayerState{
		Leader:    s.Leader,
		Wealth:    int32(s.Wealth),
		Health:    int32(s.Health),
		Stability: int32(s.Stability),
	}

	switch val := c.(type) {
	case actionCard:
		m.Card = &PlayerState_ActionCard{ActionCard: val.toProto()}
	case infoCard:
		m.Card = &PlayerState_InfoCard{InfoCard: val.toProto()}
	case votingCard:
		m.Card = &PlayerState_VotingCard{VotingCard: val.toProto()}
	default:
		return nil, fmt.Errorf("unknown card type: %v", c)
	}

	return proto.Marshal(&Message{Content: &Message_PlayerState{PlayerState: &m}})
}

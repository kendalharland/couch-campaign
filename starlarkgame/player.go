package starlarkgame

import (
	"encoding/json"
)

type PlayerState struct {
	ID string `json:"id"`

	// Card attributes.
	CardRef        CardRef `json:"card_ref" starlark:"card_ref"`
	CardSpeaker    string  `json:"card_speaker"`
	CardText       string  `json:"card_text"`
	CardAcceptText string  `json:"card_accept_text"`
	CardRejectText string  `json:"card_reject_text"`

	// Societal attributes.
	Leader             string `json:"leader"`
	LeaderTimeInOffice int    `json:"leader_time_in_office"`
	Health             int    `json:"health" starlark:"health"`
	Wealth             int    `json:"wealth" starlark:"wealth"`
	Stability          int    `json:"stability" starlark:"stability"`
}

func (s *PlayerState) ToJSONString() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (s *PlayerState) IsFailed() bool {
	return s.Wealth <= 0 || s.Health <= 0 || s.Stability <= 0
}

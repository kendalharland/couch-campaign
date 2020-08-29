package starlarkgame

import (
	"encoding/json"
)

type Card struct {
	Ref            CardRef `json:"card_ref" starlark:"ref,mutable"`
	Speaker        string  `json:"card_speaker" starlark:"speaker,mutable"`
	Text           string  `json:"card_text" starlark:"text,mutable"`
	AcceptText     string  `json:"card_accept_text" starlark:"accept_text,mutable"`
	CardRejectText string  `json:"card_reject_text" starlark:"reject_text,mutable"`
}

type PlayerState struct {
	// TODO: delete.
	CardRef            CardRef `json:"card_ref" starlark:"card_ref"`
	Leader             string  `json:"leader"`
	LeaderTimeInOffice int     `json:"leader_time_in_office"`
	Health             int     `json:"health" starlark:"health,mutable"`
	Wealth             int     `json:"wealth" starlark:"wealth,mutable"`
	Stability          int     `json:"stability" starlark:"stability,mutable"`
}

func (s *PlayerState) SetHealth(value int) {
	s.Health = value
}

func (s *PlayerState) SetWealth(value int) {
	s.Wealth = value
}
func (s *PlayerState) SetStability(value int) {
	s.Stability = value
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

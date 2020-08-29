package starlarkgame

import (
	"encoding/json"
)

type State struct {
	CardRef            CardRef `json:"card_ref" starlark:"card_ref,mutable"`
	CardSpeaker        string  `json:"card_speaker" starlark:"card_speaker,mutable"`
	CardText           string  `json:"card_text" starlark:"card_text,mutable"`
	CardAcceptText     string  `json:"card_accept_text" starlark:"card_accept_text,mutable"`
	CardRejectText     string  `json:"card_reject_text" starlark:"card_reject_text,mutable"`
	Leader             string  `json:"leader,mutable"`
	LeaderTimeInOffice int     `json:"leader_time_in_office"`
	Health             int     `json:"health" starlark:"health,mutable"`
	Wealth             int     `json:"wealth" starlark:"wealth,mutable"`
	Stability          int     `json:"stability" starlark:"stability,mutable"`
}

func (s *State) SetHealth(value int) {
	s.Health = value
}

func (s *State) SetWealth(value int) {
	s.Wealth = value
}

func (s *State) SetStability(value int) {
	s.Stability = value
}

func (s *State) SetCardRef(value string) {
	s.CardRef = CardRef(value)
}

func (s *State) SetCardSpeaker(value string) {
	s.CardSpeaker = value
}

func (s *State) SetCardText(value string) {
	s.CardText = value
}

func (s *State) SetCardAcceptText(value string) {
	s.CardAcceptText = value
}

func (s *State) SetCardRejectText(value string) {
	s.CardRejectText = value
}

func (s *State) ToJSONString() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (s *State) IsFailed() bool {
	return s.Wealth <= 0 || s.Health <= 0 || s.Stability <= 0
}

package couchcampaign

import "couchcampaign/multiplayer"

// PlayerStateView is a read-only view into a starlarkgame player's state.
type PlayerState struct {
	ID multiplayer.CID
	// Card attributes.
	CardID         string
	CardSpeaker    string
	CardText       string
	CardAcceptText string
	CardRejectText string
	// Societal attributes.
	Leader             string
	LeaderTimeInOffice int
	Health             int
	Wealth             int
	Stability          int
}

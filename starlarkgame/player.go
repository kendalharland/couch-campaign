package starlarkgame

// PlayerStateView is a read-only view into a starlarkgame player's state.
type PlayerState struct {
	ID string
	// Card attributes.
	CardRef        CardRef
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

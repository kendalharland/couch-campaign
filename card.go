package couchcampaign

// Card is a prompt/moment/action a user may take which effects their society.
type Card struct {
	ID         CardRef
	Speaker    string
	Text       string
	AcceptText string
	RejectText string
	ImageURL   string
	OnShow     []cardEffect
	OnAccept   []cardEffect
	OnReject   []cardEffect
}

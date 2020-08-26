package couchcampaign

// Input represents the input sent from the client.
type Input string

const (
	NoInput Input = ""

	// InputCardAccepted is sent when a card is accepted by the user.
	InputCardAccepted = "accept"

	// InputCardRejected is sent when a card is rejected by the user.
	InputCardRejected = "reject"
)

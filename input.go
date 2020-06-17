package couchcampaign

// Input represents the input sent from the client.
type Input string

const (
	// InputErr is sent internally when the input recieved is invalid.
	InputErr Input = ""

	// InputCardShown is sent when a card is shown to the user.
	InputCardShown = "show"

	// InputCardAccepted is sent when a card is accepted by the user.
	InputCardAccepted = "accept"

	// InputCardRejected is sent when a card is rejected by the user.
	InputCardRejected = "reject"
)

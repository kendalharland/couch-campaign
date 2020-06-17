package couchcampaign

// Input represents the input sent from the client.
type Input int

const (
	// InputErr is sent internally when the input recieved is invalid.
	InputErr Input = iota

	// InputCardShown is sent when a card is shown to the user.
	InputCardShown

	// InputCardAccepted is sent when a card is accepted by the user.
	InputCardAccepted

	// InputCardRejected is sent when a card is rejected by the user.
	InputCardRejected
)

package multiplayer

// Message represents a message sent between the server and the client.
type Message struct {
	CID  CID
	Data []byte

	// Whether to wait for a response to this message before processing the next one.
	//
	// This is only used on the server-side to send multiple messages to a client without
	// waiting for a response.
	SkipResponse bool
} 
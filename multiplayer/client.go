package multiplayer

import (
	"github.com/gorilla/websocket"
)

// CID is a type alias for a client ID.
type CID = string

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

// ClientError represents an error in server-client communication.
type ClientError struct {
	CID CID
	Err error
}

// Client is a server side thread that communicates with a remote client.
type Client struct {
	id CID
	ws *concurrentWebSocket
}

// NewClient creates a new worker for the given socket connection.
func NewClient(id CID, ws *websocket.Conn) *Client {
	return &Client{
		id: id,
		ws: newConcurrentWebsocket(ws),
	}
}

// Run executes all client updates from the input channel.
//
// It exits when either the update channel or the connection is closed.
//
// This should be run in a separate Go-routine.
func (w *Client) Run(outgoing <-chan Message, incoming chan<- Message, errs chan<- ClientError) {
	for message := range outgoing {
		input, err := w.deliver(message)
		if err != nil {
			errs <- ClientError{CID: w.id, Err: err}
			return
		}
		if input != nil {
			incoming <- Message{CID: w.id, Data: input}
		}
	}
}

func (w *Client) deliver(m Message) ([]byte, error) {
	if err := w.send(m); err != nil {
		return nil, err
	}
	if m.SkipResponse {
		return nil, nil
	}
	return w.recieve()
}

// Sends the update to the remote peer.
func (w *Client) send(message Message) error {
	return w.ws.Write(message.Data)
}

// Receives a messge from the remote peer, which may be nil
func (w *Client) recieve() ([]byte, error) {
	_, input, err := w.ws.Read()
	return input, err
}

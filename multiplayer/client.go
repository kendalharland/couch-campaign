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
func (w *Client) Run(updates <-chan Message, messages chan<- Message, errs chan<- ClientError) {
	for update := range updates {
		input, err := w.deliver(update)
		if err != nil {
			errs <- ClientError{
				CID: w.id,
				Err: err,
			}
		} else {
			messages <- Message{
				CID:  w.id,
				Data: input,
			}
		}
	}
}

// Deliver sends the update to the remote client.
//
// Returns the remote client's response, or nil if the update does not require
// a response.
func (w *Client) deliver(message Message) ([]byte, error) {
	if err := w.ws.Write(message.Data); err != nil {
		return nil, err
	}
	_, input, err := w.ws.Read()
	return input, err
}

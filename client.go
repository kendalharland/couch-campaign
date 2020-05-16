package couchcampaign

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ClientJob is a unit of work for the ClientDriver.
type ClientJob struct {
	Card  Card
	Stats playerState
}

// ClientMessage represents game input from the client.
type ClientMessage struct {
	PID  PID
	Type ClientMessageType

	Card  Card
	Input ClientInput
}

// ClientError represents an error in server-client communication.
type ClientError struct {
	PID PID
	err error
}

func (e ClientError) Error() string {
	return e.err.Error()
}

// ClientMessageType describes the kind of message recieved from the client.
type ClientMessageType int

// ClientMessage type constants.
const (
	// InputClientMessage is used when the client's message contains game input.
	//
	// This is sent when the player plays a card.
	InputClientMessage ClientMessageType = iota

	// DisconnectClientMessage is used when the client has disconnected from the server.
	DisconnectClientMessage
)

type ClientInput int

// ClientInput constants.
const (
	NoInput ClientInput = iota
	DismissInfoCardInput
	AcceptActionCardInput
	RejectActionCardInput
)

// ClientWorker implements a server-side go-routine that communicates with a single client.
type ClientWorker struct {
	PID    PID
	driver *clientDriver
}

// NewClientWorker creates a new worker for the given socket connection.
func NewClientWorker(pid PID, ws *websocket.Conn) *ClientWorker {
	return &ClientWorker{
		PID:    pid,
		driver: &clientDriver{newConcurrentWebsocket(ws)},
	}
}

// Run executes each job from the input channel.
//
// It exits either when all jobs are complete or when the connection is closed.
//
// This should be run in a separate Go-routine.
func (w *ClientWorker) Run(jobs <-chan ClientJob, messages chan<- ClientMessage, errs chan<- ClientError) {
	for job := range jobs {
		input, err := w.do(job)
		if err != nil {
			errs <- ClientError{
				PID: w.PID,
				err: err,
			}
		} else if input != NoInput {
			messages <- ClientMessage{
				PID:   w.PID,
				Type:  InputClientMessage,
				Card:  job.Card,
				Input: input,
			}
		}
	}
}

// do does the given ClientJob.
//
// Returns the player's input action, or the empty string if the job did not
// not require any input.
func (w *ClientWorker) do(job ClientJob) (ClientInput, error) {
	if err := w.driver.sendState(job.Card, job.Stats); err != nil {
		return -1, err
	}
	if !cardRequiresInput(job.Card) {
		return -1, nil

	}
	input, err := w.driver.getInput()
	if err != nil {
		return -1, err
	}

	switch input {
	case "accept":
		if _, ok := job.Card.(actionCard); ok {
			return AcceptActionCardInput, nil
		}
		return DismissInfoCardInput, nil
	case "reject":
		if _, ok := job.Card.(actionCard); ok {
			return RejectActionCardInput, nil
		}
		return DismissInfoCardInput, nil
	default:
		return NoInput, fmt.Errorf("invalid input: %q", input)
	}
}

type clientDriver struct {
	ws *concurrentWebSocket
}

func (d *clientDriver) close() error {
	return d.ws.Close()
}

func (d *clientDriver) getInput() (string, error) {
	_, input, err := d.ws.Read()
	if err != nil {
		return "", err
	}
	return string(input), nil
}

func (d *clientDriver) sendState(c Card, s playerState) error {
	m, err := renderCard(c, s)
	if err != nil {
		return fmt.Errorf("renderCard: %v", err)
	}
	if err := d.ws.Write(m); err != nil {
		return err
	}
	return nil
}

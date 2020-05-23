package cardgame

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ClientJob is a unit of work for the ClientDriver.
type ClientJob struct {
	PID           PID
	State         []byte
	RequiresInput bool
}

// ClientMessage represents game input from the client.
type ClientMessage struct {
	PID   PID
	Card  CardRef
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

// ClientInput is a signal from the client to the game server.
type ClientInput int

// ClientInput constants.
const (
	// FailedInput is sent when an error occurred either displaying a card or receiving input.
	FailedInput ClientInput = iota

	// NoINput is sent when a card was displayed but input was not required.
	NoInput

	// AcceptCardInput is sent when the client accepts a card.
	AcceptCardInput

	// AcceptCardInput is sent when the client rejects a card.
	RejectCardInput
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
	if err := w.driver.sendState(job.State); err != nil {
		return FailedInput, err
	}
	if !job.RequiresInput {
		return NoInput, nil
	}
	input, err := w.driver.getInput()
	if err != nil {
		return FailedInput, err
	}

	// TODO: Replace this in favor of the client sending the input code directly.
	switch input {
	case "accept":
		return AcceptCardInput, nil
	case "reject":
		return RejectCardInput, nil
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

func (d *clientDriver) sendState(state []byte) error {
	if err := d.ws.Write(state); err != nil {
		return err
	}
	return nil
}

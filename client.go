package couchcampaign

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// ClientJob is a unit of work for the ClientDriver.
type ClientJob struct {
	Card  Card
	Stats stats
}

// ClientMessage represents game input from the client.
type ClientMessage struct {
	pid  PID
	Type ClientMessageType

	card  Card
	input string
}

// ClientMessageType describes the kind of message recieved from the client.
type ClientMessageType int

const (
	// InputClientMessage is used when the client's message contains game input.
	//
	// This is sent when the player plays a card.
	InputClientMessage ClientMessageType = iota

	// DisconnectClientMessage is used when the client has disconnected from the server.
	DisconnectClientMessage
)

// ClientDriver implements a server-side go-routine that communicates with a single client.
type ClientDriver struct {
	pid  PID
	conn *connection
	out  chan<- ClientMessage
}

func NewClientDriverFromWebSocket(pid PID, c *websocket.Conn) *ClientDriver {
	return NewClientDriver(newConnection(c), pid)
}

func NewClientDriver(c *connection, pid PID) *ClientDriver {
	return &ClientDriver{
		conn: c,
		pid:  pid,
	}
}

func (cli *ClientDriver) setOutChan(out chan<- ClientMessage) {
	cli.out = out
}

// Run executes each job from the input channel.
//
// It exits either when all jobs are complete or when the connection is closed.
//
// This should be run in a separate Go-routine.
func (cli *ClientDriver) Run(jobs <-chan ClientJob) {
	for job := range jobs {
		if err := cli.do(job); err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr) {
				cli.out <- ClientMessage{
					pid:  cli.pid,
					Type: DisconnectClientMessage,
				}
				return
			}
			log.Printf("ClientDriver %v do: %v", cli.pid, err)
		}
	}
}

func (cli *ClientDriver) do(job ClientJob) error {
	if err := cli.showCard(job.Card, job.Stats); err != nil {
		return err
	}

	if cardRequiresInput(job.Card) {
		input, err := cli.getInput()
		if err != nil {
			return err
		}
		cli.out <- ClientMessage{
			pid:   cli.pid,
			Type:  InputClientMessage,
			card:  job.Card,
			input: input,
		}
		return nil
	}
	return nil
}

func (cli *ClientDriver) close() error {
	return cli.conn.Close()
}

func (cli *ClientDriver) getInput() (string, error) {
	_, input, err := cli.conn.Read()
	if err != nil {
		return "", err
	}
	return string(input), nil
}

func (cli *ClientDriver) showCard(c Card, s stats) error {
	m, err := renderCard(c, s)
	if err != nil {
		return fmt.Errorf("renderCard: %v", err)
	}
	if err := cli.conn.WriteBinary(m); err != nil {
		return err
	}
	return nil
}

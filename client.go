package couchcampaign

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	pid  PID
	conn *connection
}

func NewClientFromWebSocket(pid PID, c *websocket.Conn) *Client {
	conn := newConnection(c)
	return newClient(conn, pid)
}

func newClient(c *connection, pid PID) *Client {
	return &Client{
		conn: c,
		pid:  pid,
	}
}

func (cli *Client) run(jobs <-chan func(*Client)) {
	for job := range jobs {
		job(cli)
	}
}

func (cli *Client) getInput() (string, error) {
	_, message, err := cli.conn.Read()
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (cli *Client) showCard(c Card, s *stats) {
	m, err := renderCard(c, s)
	if err != nil {
		log.Printf("renderCard: %v", err)
	}
	cli.conn.WriteBinary(m)
}

func (cli *Client) Close() error {
	return cli.conn.Close()
}

package couchcampaign

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// A websocket connection safe for concurrent use.
type connection struct {
	c  *websocket.Conn
	mu *sync.RWMutex
}

func newConnection(c *websocket.Conn) *connection {
	return &connection{c: c, mu: &sync.RWMutex{}}
}

func (c *connection) Read() (messageType int, p []byte, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.c.ReadMessage()
}

func (c *connection) Write(message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	w, err := c.c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err1 := fmt.Fprintf(w, message)
	err2 := w.Close()
	if err1 != nil {
		return err
	}
	return err2
}

func (c *connection) WriteBinary(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	w, err := c.c.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	_, err1 := w.Write(data)
	err2 := w.Close()
	if err1 != nil {
		return err
	}
	return err2
}

func (c *connection) Close() error {
	return c.c.Close()
}

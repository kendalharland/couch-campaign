package cardgame

import (
	"sync"

	"github.com/gorilla/websocket"
)

func IsConnectionCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr)
}

type concurrentWebSocket struct {
	c  *websocket.Conn
	mu *sync.RWMutex
}

func newConcurrentWebsocket(c *websocket.Conn) *concurrentWebSocket {
	return &concurrentWebSocket{c: c, mu: &sync.RWMutex{}}
}

func (c *concurrentWebSocket) Read() (messageType int, p []byte, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.c.ReadMessage()
}

func (c *concurrentWebSocket) Write(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	w, err := c.c.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	_, err1 := w.Write(data)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

func (c *concurrentWebSocket) Close() error {
	return c.c.Close()
}

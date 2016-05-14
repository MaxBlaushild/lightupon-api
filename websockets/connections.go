package websockets

import (
	"github.com/gorilla/websocket"
	"lightupon-api/models"
	"log"
	"time"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Content []byte
	Passcode string
}

type Connection struct {
	WS *websocket.Conn
	Passcode string
	Send chan models.PullResponse
	User models.User
}

func (c *Connection) ReadPump() {
	defer func() {
		H.Unregister <- c
		c.WS.Close()
	}()

	c.WS.SetReadLimit(maxMessageSize)
	c.WS.SetReadDeadline(time.Now().Add(pongWait))
	c.WS.SetPongHandler(func(string) error { 
		c.WS.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})

	for {
		_, message, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		RouteIncomingMessage(message, c)
	}
}

func (c *Connection) Write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	message := Message{Content: payload}
	return c.WS.WriteJSON(message)
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case pullResponse, ok := <- c.Send:

			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.WS.WriteJSON(pullResponse); err != nil {
				return
			}

		case <-ticker.C:

			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}

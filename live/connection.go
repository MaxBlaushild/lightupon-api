package live

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
	"encoding/json"
	"bytes"
  "fmt"
)

type Response struct {
	NextSceneAvailable bool
	CurrentSceneOrder int
	Passcode string
	Event string
	UpdatedPostID uint
}

type PostUpdate struct {
	UpdatedPostID uint
	UserID uint
}

type Location struct {
	Latitude float64
	Longitude float64
}

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 1) / 10
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
	Send chan Response
	UserID uint
	Location Location
}

func (c *Connection) ReadPump() {
	defer func() {
		fmt.Print("The read pump died a natural death.")
		Hub.Unregister <- c
		c.WS.Close()
	}()

	c.ConfigureRead()

	for {
		_, locationBytes, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.SetLocation(locationBytes)
		c.UpdateClient()
	}
}

func (c *Connection) ConfigureRead() {
	c.WS.SetReadLimit(maxMessageSize)
	c.WS.SetReadDeadline(time.Now().Add(pongWait))
	c.WS.SetPongHandler(func(string) error { 
		c.WS.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})
}

func (c *Connection) SetLocation(locationBytes []byte) {
	location := Location{}
	buffer := bytes.NewBuffer(locationBytes)
	decoder := json.NewDecoder(buffer)
	decoder.Decode(&location)
	c.Location = location
}

func (c *Connection) Write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	return c.WS.WriteMessage(mt, payload)
}

func (c *Connection) UpdateClient() {
	response := Response{ Passcode: c.Passcode }
  Hub.Broadcast <- response
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()

	for {
		select {
		case response, ok := <- c.Send:
			if !ok {
				Hub.Unregister <- c
				return
			}
			if err := c.WS.WriteJSON(response); err != nil {
				return
			}
		case <- ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}
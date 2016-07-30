package websockets

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
	"encoding/json"
	"lightupon-api/models"
	"bytes"
    "fmt"
    "github.com/davecgh/go-spew/spew"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Hour
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

	c.ConfigureRead()

	for {
		_, message, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.UpdateLocation(message)
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

func (c *Connection) UpdateLocation(message []byte) {
	fmt.Println("stand next to this money like HEY HEY")
	spew.Dump(location)
	location := models.Location{}
	buffer := bytes.NewBuffer(message)
  decoder := json.NewDecoder(buffer)
  err := decoder.Decode(&location); if err != nil {
  	fmt.Println(err)
  }
	c.User.Location = location
}

func (c *Connection) Write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	return c.WS.WriteMessage(mt, payload)
}

func (c *Connection) UpdateClient() {
	party := models.Party{}
	models.DB.Preload("Scene.Cards").Where("passcode = ?", c.Passcode).First(&party)
  H.Broadcast <- party
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
				H.Unregister <- c
				return
			}
			if err := c.WS.WriteJSON(pullResponse); err != nil {
				return
			}
		case <- ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}

		}
	}
}
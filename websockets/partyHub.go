package websockets

import (
	"fmt"
	"lightupon-api/models"
)

type hub struct {
	Connections map[string]map[*Connection]bool
	Broadcast chan models.PullResponse
	Register chan *Connection
	Unregister chan *Connection
}

var H = hub{
	Broadcast:   make(chan models.PullResponse),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	Connections: make(map[string]map[*Connection]bool),
}

func (h *hub) StartHub() {
	fmt.Println("[partyHub] ready for connections")
	for {
		select {
		case c := <- h.Register:
			h.RegisterConnection(c)
		case c := <-h.Unregister:
			h.UnregisterConnection(c)
		case pullResponse := <- h.Broadcast:
			h.GatherUsers(&pullResponse)
			h.PushToParty(pullResponse)
		}
	}
}

func (h *hub) RegisterConnection(c *Connection) {
	h.Connections[c.Passcode] = make(map[*Connection]bool)
	h.Connections[c.Passcode][c] = true
}

func (h *hub) GatherUsers(pullResponse *models.PullResponse) {
	for c := range h.Connections[pullResponse.Passcode] {
		pullResponse.Users = append(pullResponse.Users, c.User)
	}
}

func (h *hub) PushToParty(pullResponse models.PullResponse) {
	for c := range h.Connections[pullResponse.Passcode] {
		select {
		case c.Send <- pullResponse:
		default:
			close(c.Send)
			delete(h.Connections[c.Passcode], c)
		}
	}
}

func (h *hub) UnregisterConnection(c *Connection) {
	if _, ok := h.Connections[c.Passcode][c]; ok {
		delete(h.Connections[c.Passcode], c)
		close(c.Send)
	}
}

func (h *hub) DeactivateUser(user models.User, passcode string) {
	for c := range h.Connections[passcode] {
		if user.ID == c.User.ID {
			h.Unregister <- c
		}
	}
}
package websockets

import (
	"fmt"
	"lightupon-api/models"
    "github.com/davecgh/go-spew/spew"
)

type hub struct {
	Connections map[string]map[*Connection]bool
	Broadcast chan models.Party
	Register chan *Connection
	Unregister chan *Connection
}

var H = hub{
	Broadcast:   make(chan models.Party),
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
		case party := <- h.Broadcast:
			pullResponse := h.CreatePullResponse(party)
			h.PushToParty(pullResponse)
		}
	}
}

func (h *hub) RegisterConnection(c *Connection) {
	h.Connections[c.Passcode] = make(map[*Connection]bool)
	h.Connections[c.Passcode][c] = true
}

func (h *hub) CreatePullResponse(party models.Party) models.PullResponse {
  pullResponse := models.PullResponse{Passcode: party.Passcode, Scene: party.Scene, NextScene: party.NextScene()}
  pullResponse.Users = h.GatherUsersFromParty(party)
  pullResponse.NextSceneAvailable = h.IsNextSceneAvailable(party)
  fmt.Print("pullResponse")
  spew.Dump(pullResponse)
  return pullResponse
}

func (h *hub) IsNextSceneAvailable(party models.Party)(nextSceneAvailable bool) {
	nextScene := party.NextScene()
	for c := range h.Connections[party.Passcode] {
		nextSceneAvailable = nextSceneAvailable || c.User.IsAtScene(nextScene)
	}
	return
}

func (h *hub) GatherUsersFromParty(party models.Party)(users []models.User) {
	for c := range h.Connections[party.Passcode] {
		users = append(users, c.User)
	}
	return
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
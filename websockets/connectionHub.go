package websockets

import (
	"fmt"
	"lightupon-api/models"
	"github.com/davecgh/go-spew/spew"
)

type hub struct {
	PartyConnections map[string]map[*Connection]bool
	Connections map[string]*Connection
	Broadcast chan models.Party
	Register chan *Connection
	Unregister chan *Connection

}

var H = hub{
	Broadcast:   make(chan models.Party),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	PartyConnections: make(map[string]map[*Connection]bool),
	Connections: make(map[string]*Connection),
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

func (h *hub) RegisterPartyConnection(c *Connection) {
	if (h.PartyConnections[c.Passcode] == nil) {
		h.PartyConnections[c.Passcode] = make(map[*Connection]bool)
	}
	h.PartyConnections[c.Passcode][c] = true
	fmt.Println("Registered party connection")
}

func (h *hub) RegisterConnection(c *Connection) {
	h.Connections[c.User.FacebookId] = c
	fmt.Println("Registered connection")
	if (len(c.Passcode) > 0) {
		h.RegisterPartyConnection(c)
	}
}

func (h *hub) AddUserConnectionToParty(user models.User, party models.Party) {
	c := h.Connections[user.FacebookId]
	c.Passcode = party.Passcode
	h.RegisterPartyConnection(c)
}

func (h *hub) CreatePullResponse(party models.Party) models.PullResponse {
  pullResponse := models.PullResponse{Passcode: party.Passcode, Scene: party.Scene, NextScene: party.NextScene()}
  pullResponse.Users = h.GatherUsersFromParty(party)
  pullResponse.NextSceneAvailable = h.IsNextSceneAvailable(party)
  return pullResponse
}

func (h *hub) IsNextSceneAvailable(party models.Party)(nextSceneAvailable bool) {
	nextScene := party.NextScene()
	for c := range h.PartyConnections[party.Passcode] {
		nextSceneAvailable = nextSceneAvailable || c.User.IsAtScene(nextScene)
	}
	return
}

func (h *hub) GatherUsersFromParty(party models.Party)(users []models.User) {
	for c := range h.PartyConnections[party.Passcode] {
		users = append(users, c.User)
	}
	return
}

func (h *hub) PushToParty(pullResponse models.PullResponse) {
	if len(pullResponse.Passcode) > 0 {
		fmt.Println(pullResponse.Passcode)
	}
	for c := range h.PartyConnections[pullResponse.Passcode] {
		fmt.Println("pushing this ish to connection" )
		spew.Dump(h.PartyConnections[pullResponse.Passcode])
		select {
		case c.Send <- pullResponse:
			fmt.Println(c.User.FacebookId)
		default:
			fmt.Println("Pull request did not get broadcast correctly.")
			H.UnregisterConnection(c)
		}
	}
}

func (h *hub) UnregisterConnection(c *Connection) {
	if _, ok := h.PartyConnections[c.Passcode][c]; ok {
		delete(h.PartyConnections[c.Passcode], c)
		close(c.Send)
	}
	delete(h.Connections, c.User.FacebookId)
}

func (h *hub) DeactivateUserFromParty(user models.User, passcode string) {
	c := h.Connections[user.FacebookId]
	fmt.Println("deleting connection")
	fmt.Println(c.Passcode)
	if _, ok := h.PartyConnections[c.Passcode][c]; ok {
		fmt.Println("found that sucker")
		delete(h.PartyConnections[c.Passcode], c)
		c.Passcode = ""
	}
}
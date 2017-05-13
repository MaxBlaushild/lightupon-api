package live

import(
	"fmt"
)

type hub struct {
	Parties map[string]Party
	Connections map[uint]*Connection
	Broadcast chan Response
	UpdateClient chan SceneUpdate
	Register chan *Connection
	Unregister chan *Connection
	PutParty chan Party
	EndParty chan string
}

func (h *hub) Start() {

	fmt.Println("we are [live]")
	for {
		select {
		case c := <- h.Register:
			h.RegisterConnection(c)
		case party := <- h.PutParty:
			h.SyncParty(party)
		case sceneUpdate := <- h.UpdateClient:
			response := Response{UpdatedSceneID: sceneUpdate.UpdatedSceneID}
			if h.Connections[sceneUpdate.UserID] != nil {
				h.Connections[sceneUpdate.UserID].Send <- response
			}
		case passcode := <- h.EndParty:
			h.UnregisterParty(passcode)
		case c := <- h.Unregister:
			h.UnregisterConnection(c)
		case response := <- h.Broadcast:
			party := h.Parties[response.Passcode]

			if (party.Exists()) {
				party.Push(response)
			}
		}
	}
}

func (h *hub) SyncParty(party Party) {
	existingParty := h.Parties[party.Passcode]

	if (existingParty.Users != nil) {
		party.Users = existingParty.Users
	}

	h.Parties[party.Passcode] = party
}

func (h *hub) UnregisterParty(passcode string) {
	delete(h.Parties, passcode)
}

func (h *hub) RegisterPartyConnection(c *Connection) {
	party := h.Parties[c.Passcode]

	if (!party.Exists()) {
		party.Users = make(map[uint]*Connection)
	}

	party.Users[c.UserID] = c
}

func (h *hub) RegisterConnection(c *Connection) {
	h.Connections[c.UserID] = c
	if (len(c.Passcode) > 0) {
		h.RegisterPartyConnection(c)
	}
}

func (h *hub) AddUserToParty(userID uint, passcode string) {
	c := h.Connections[userID]
	c.Passcode = passcode
	h.RegisterPartyConnection(c)
}

func (h *hub) UnregisterConnection(c *Connection) {
	if _, ok := h.Parties[c.Passcode].Users[c.UserID]; ok {
		delete(h.Parties[c.Passcode].Users, c.UserID)
		close(c.Send)
	}
	delete(h.Connections, c.UserID)
}

func (h *hub) DropUserFromParty(userID uint, passcode string) {
	c := h.Connections[userID]; if c != nil {
		party := h.Parties[passcode]
		found := party.Users[userID]; if found != nil {
			delete(party.Users, userID)
			c.Passcode = ""
		}
	}

}
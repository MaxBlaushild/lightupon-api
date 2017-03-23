package live

var Hub = hub{
	Broadcast:   make(chan Response),
	Connections: make(map[uint]*Connection),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	Parties: make(map[string]Party),
	EndParty: make(chan string),
	PutParty: make(chan Party),
}
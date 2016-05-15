package models

type PullResponse struct {
	AdvanceToNextScene bool
	NextScene Scene
	Passcode string
	Action string
	Users []User
}
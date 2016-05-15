package models

type PullResponse struct {
	NextSceneAvailable bool
	NextScene Scene
	Passcode string
	Action string
	Users []User
}
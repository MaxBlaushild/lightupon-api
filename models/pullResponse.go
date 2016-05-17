package models

type PullResponse struct {
	NextSceneAvailable bool
	Scene Scene
	Passcode string
	Action string
	Users []User
}
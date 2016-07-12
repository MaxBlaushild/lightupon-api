package models

type PullResponse struct {
	NextSceneAvailable bool
	Scene Scene
	NextScene Scene
	Passcode string
	Action string
	Users []User
}
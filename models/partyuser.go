package models

import (
	"github.com/jinzhu/gorm"
)

type Partyuser struct {
	gorm.Model
	User User
	UserID uint `gorm:"index"`
	Party Party
	PartyID uint `gorm:"index"`
	CurrentSceneOrderID uint `gorm:"default:0"`
	Scene Scene
	Completed bool `gorm:"default:false"`
}

func (p *Partyuser) IsUserAtNextScene(lat float64, lon float64) (isAtNextScene bool, nextScene Scene) {
  DB.Where("trip_id = ? AND scene_order = ?", p.Party.TripID, p.CurrentSceneOrderID + 1).First(&nextScene)
  distanceFromScene := CalculateDistance(UserLocation{Latitude: lat, Longitude: lon}, UserLocation{Latitude: nextScene.Latitude, Longitude: nextScene.Latitude})
  isAtNextScene = distanceFromScene < unlockThresholdSmall
  return
}
package models

import (
	"github.com/jinzhu/gorm"
	"math"
)

type Partyuser struct {
	gorm.Model
	User User
	UserID uint `gorm:"index"`
	Party Party
	PartyID uint `gorm:"index"`
	CurrentSceneOrderID uint `gorm:"default:0"`
	Scene Scene
}

func (p *Partyuser) IsUserAtNextScene(lat float64, lon float64) (isAtNextScene bool, nextScene Scene) {
  DB.Where("trip_id = ? AND scene_order = ?", p.Party.TripID, p.CurrentSceneOrderID + 1).First(&nextScene)

  // Decide whether we're at the next scene
  latDiff := nextScene.Latitude - lat
  lonDiff := nextScene.Longitude - lon
  distanceFromScene := math.Pow(latDiff, 2) + math.Pow(lonDiff, 2)
  isAtNextScene = distanceFromScene < threshold

  return
}
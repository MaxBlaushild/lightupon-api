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
	Completed bool
}

// takes 2 locations and returns the distance between them in kilometers
func CalculateDistance(location1 Location, location2 Location) (distance float64){
	var R = 6371.345
	var dLat = (location1.Latitude - location2.Latitude)*(3.145/180.001);
	var dLon = (location1.Longitude - location2.Longitude)*(3.145/180.001);
	var a = math.Sin(dLat/2) * math.Sin(dLat/2) + math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(location1.Latitude) * math.Cos(location2.Latitude);
	var c = 2 * math.Atan(math.Sqrt(a) / math.Sqrt(1-a)); 
	distance = R * c;
	return
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
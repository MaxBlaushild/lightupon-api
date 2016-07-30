package models

import (
	"github.com/jinzhu/gorm"
	"math"
	"fmt"
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

func CalculateDistance(location1 Location, location2 Location) (distance float64){
	// dlon := location1.Longitude - location2.Longitude
	dlat := location1.Latitude - location2.Latitude
	// distance := dlat

	fmt.Println(dlat)
	return
	// a := math.Pow((sin(dlat/2)),2) + cos(lat1) * cos(lat2) * math.Pow(sin(dlon/2),2) 
// c := 2 * atan2( sqrt(a), sqrt(1-a) ) 
// d := R * c (where R is the radius of the Earth)

// a := sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
// c = 2 ⋅ atan2( √a, √(1−a) )
// d = R ⋅ c

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
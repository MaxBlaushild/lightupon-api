package models

import(
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/aws"
      "math"
      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint `gorm:"index"`
  BackgroundUrl string `gorm:"not null"`
  SceneOrder uint `gorm:"not null"`
  Featured bool
  Cards []Card
  SoundKey string
  SoundResource string
}

func ShiftScenesUp(sceneOrder int, tripID int) bool {
  scene := Scene{}
  DB.Where("trip_id = $1 AND scene_order = $2", tripID, sceneOrder).First(&scene)
  if scene.ID == 0 {
    return true
  } else {
    ShiftScenesUp(sceneOrder + 1, 1)
    DB.Model(&scene).Update("scene_order", sceneOrder + 1)
    return true
  }
}

func ShiftScenesDown(sceneOrder int, tripID int) bool {
  scene := Scene{}
  DB.Where("trip_id = $1 AND scene_order = $2", tripID, sceneOrder + 1).First(&scene)
  if scene.ID == 0 {
    return true
  } else {
    ShiftScenesDown(sceneOrder + 1, 1)
    DB.Model(&scene).Update("scene_order", sceneOrder)
    return true
  }
}

func (s *Scene) IsAtScene(location UserLocation)(isAtNextScene bool) {
  sceneLocation := UserLocation{Latitude: s.Latitude, Longitude: s.Longitude}
  distanceFromScene := CalculateDistance(location, sceneLocation)
  isAtNextScene = distanceFromScene < threshold
  return
}

// takes 2 locations and returns the distance between them in kilometers
func CalculateDistance(location1 UserLocation, location2 UserLocation) (distance float64){
  var R = 6371.345
  var dLat = (location1.Latitude - location2.Latitude)*(3.145/180.001);
  var dLon = (location1.Longitude - location2.Longitude)*(3.145/180.001);
  var a = math.Sin(dLat/2) * math.Sin(dLat/2) + math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(location1.Latitude) * math.Cos(location2.Latitude);
  var c = 2 * math.Atan(math.Sqrt(a) / math.Sqrt(1-a)); 
  distance = R * c;
  return
}

func (s *Scene) PopulateSound() {
  url, err := aws.GetAsset("audio", s.SoundKey)

  if err != nil {
    fmt.Println(err)
  }

  s.SoundResource = url
}
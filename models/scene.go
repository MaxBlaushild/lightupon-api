package models

import(
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/aws"
      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint
  BackgroundUrl string
  SceneOrder uint
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

func (s *Scene) PopulateSound() {
  url, err := aws.GetAsset("audio", s.SoundKey)

  if err != nil {
    fmt.Println(err)
  }

  s.SoundResource = url
}
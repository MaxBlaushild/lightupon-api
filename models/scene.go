package models

import(
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/services/aws"
      "lightupon-api/services/googleMaps"

      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint `gorm:"index"`
  Trip Trip
  BackgroundUrl string `gorm:"not null"`
  SceneOrder uint `gorm:"not null"`
  Cards []Card
  Comments []Comment
  SceneLikes []SceneLike 
  GooglePlaceID string
  Route string
  User User
  UserID uint
  FormattedAddress string
  Locality string
  Neighborhood string
  PostalCode string
  Country string
  AdministrativeLevelTwo string
  AdministrativeLevelOne string
  StreetNumber string
  SoundKey string
  SoundResource string
  ConstellationPoint ConstellationPoint
  Liked bool `sql:"-"`
}

func (s *Scene) BerforeCreate() {
  place := googleMaps.GetPrettyPlace(s.Latitude, s.Longitude)
  s.FormattedAddress = place["FormattedAddress"]
  s.StreetNumber = place["street_number"]
  s.Route = place["route"]
  s.Neighborhood = place["neighborhood"]
  s.Locality = place["locality"]
  s.AdministrativeLevelTwo = place["administrative_area_level_2"]
  s.AdministrativeLevelOne = place["administrative_area_level_1"]
  s.Country = place["country"]
  s.PostalCode = place["postal_code"]
  s.GooglePlaceID = place["PlaceID"]
}

func (s *Scene) UserHasLiked(u *User) (userHasLiked bool) {
  for _, like := range s.SceneLikes {
    if like.UserID == u.ID {
      userHasLiked = true
    }
  }
  return
}

func IndexScenes() (scenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Order("created_at desc").Find(&scenes)
  return
}

func GetScenesForUser(userID string) (scenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Order("created_at desc").Where("user_id = ?", userID).Find(&scenes)
  return
}

func (s *Scene) AppendCard(card Card) (err error) {
  cardOrder := uint(len(s.Cards) + 1)
  card.CardOrder = cardOrder
  card.SceneID = s.ID
  err = DB.Save(&card).Error
  return
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

func (s *Scene) PopulateSound() {
  url, err := aws.GetAsset("audio", s.SoundKey)

  if err != nil {
    fmt.Println(err)
  }

  s.SoundResource = url
}
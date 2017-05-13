package models

import(
      "strconv"
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/services/imageMagick"
      "io/ioutil"
      "lightupon-api/services/aws"
      "net/http"
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
  PinUrl string
  SelectedPinUrl string
  ConstellationPoint ConstellationPoint
  Liked bool `sql:"-"`
  Hidden bool `sql:"-"`
  Blur float64 `sql:"-"`
}


func (s *Scene) AfterCreate(tx *gorm.DB) (err error) {
  err = s.SetPins()
  err = tx.Save(s).Error
  return
}

func (s *Scene) UserHasLiked(u *User) (userHasLiked bool) {
  for _, like := range s.SceneLikes {
    if like.UserID == u.ID {
      userHasLiked = true
    }
  }
  return
}

func GetSceneByID(sceneID string) (scene Scene, err error) {
  err = DB.Preload("Trip.User").Preload("Cards").Where("id = ?", sceneID).Find(&scene).Error
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

func (s *Scene) AppendCard(card *Card) (err error) {
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

func (s *Scene) DownloadImage() (imageBinary []byte){
  resp, err := http.Get(s.BackgroundUrl)

  defer resp.Body.Close()

  imageBinary, err = ioutil.ReadAll(resp.Body); if err != nil {
    fmt.Println("ioutil.ReadAll -> %v", err)
  }
  return 
}

func (s *Scene) SetPins() (err error) {
  imageBinary := s.DownloadImage()
  pinBinary := imageMagick.CropPin(imageBinary)
  s.PinUrl, err = s.uploadPin(pinBinary, "pin")
  return
}

func (s *Scene) getAssetName(name string) string {
  return "/scenes/" + fmt.Sprint(s.ID) + "/" + name
}

func (s *Scene) uploadPin(binary []byte, name string) (getUrl string, err error){
  asset := aws.Asset{
    Type: "images", 
    Name: s.getAssetName(name), 
    Extension: ".png",
    Binary: binary,
  }
  getUrl, err = aws.UploadAsset(asset)
  return
}

func GetFollowingScenesNearLocation(lat string, lon string, userID uint) (scenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Order("((scenes.latitude - " + lat + ")^2.0 + ((scenes.longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc").Limit(20).Find(&scenes)
  return
}

func GetScenesNearLocation(lat string, lon string, userID uint, radius string, numScenes int) (scenes []Scene, err error) {
  distanceString :="((scenes.latitude - " + lat + ")^2.0 + ((scenes.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Where(distanceString + " < " + radius).Order(distanceString + " asc").Limit(numScenes).Find(&scenes)
  for i, _ := range scenes {
    scenes[i].GetExposure(userID)
  }
  return
}

func (s *Scene) GetExposure(userID uint) (err error) {
  exposedScene := ExposedScene{}
  err = DB.Where("user_id = ? AND scene_id = ?", userID, s.ID).First(&exposedScene).Error; if err != nil {
    s.Hidden = !exposedScene.Unlocked
    s.Blur = exposedScene.Blur
  }
  return
}

// // There's not a lot of awesome abstraction here, but this could get computationally expensive so I'm trying to optimize for speed. It also requires some splainin so read the comments.
// func (scene *Scene) discover(userID uint, userLat string, userLon string) {
//   // Try to get the current blur level
//   oldExposedScene := ExposedScene{UserID : userID, SceneID : scene.ID}
//   DB.First(&oldExposedScene, oldExposedScene)

//   if (oldExposedScene.Unlocked) { // UU, UB, UL
//     // If we found a record and it's fully unlocked, then we're done. Just set the properties - no need to persist anything.
//     scene.Blur = 0.0
//     scene.Hidden = false
//   } else { // BU, BB, BL, LU, LB, LL
//     userLatFloat, _ := strconv.ParseFloat(userLat, 64)
//     userLonFloat, _ := strconv.ParseFloat(userLon, 64)
//     distanceFromScene := CalculateDistance(UserLocation{Latitude: userLatFloat, Longitude: userLonFloat}, UserLocation{Latitude: scene.Latitude, Longitude: scene.Longitude})
//     newBlur := calculateBlur(distanceFromScene)
//     if (distanceFromScene < unlockThresholdSmall) { // LU, BU
//       // Moving to Unlocked from unlockThresholdLarge non-Unlocked state, so we need to return and persist the new shit
//       scene.Blur = 0.0
//       scene.Hidden = false
//       oldExposedScene.upsertExposedScene(newBlur, scene.ID, userID, false)
//     } else {
//       if (newBlur > oldExposedScene.Blur) { // MM(change), LM
//         // save the new blur and return that new shit
//         scene.Blur = newBlur
//         scene.Hidden = true
//         oldExposedScene.upsertExposedScene(newBlur, scene.ID, userID, true)
//       } else { // MM(static), ML, LL
//         // save nothing and return the old shit
//         scene.Blur = oldExposedScene.Blur
//         scene.Hidden = true
//       }
//     }
//   }
// }

func calculateBlur(distance float64) (blur float64) {
  if (distance < unlockThresholdSmall) {
    blur = 0.0
  } else if (distance > unlockThresholdLarge) {
    blur = 1.0
  } else {
    // TODO: Update this to be a nice smoove cosine function
    blur = (distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall)
  }
  return
}

// func possiblyRecomputeAllDiscovery(lat string, lon string, userID uint) {
//   exposedScene := ExposedScene{}
//   DB.First(&exposedScene)
//   if exposedScene.ID == 0 {
//     recomputeAllDiscovery()
//   }
// }

// func recomputeAllDiscovery() {
//   fmt.Println("NOTICE: Recomputing all discovery!")
//   locations := []Location{}
//   DB.Find(&locations)
//   for i := 0; i < len(locations); i++ {
//     lat := strconv.FormatFloat(locations[i].Latitude, 'E', -1, 64)
//     lon := strconv.FormatFloat(locations[i].Longitude, 'E', -1, 64)
//     scenes := []Scene{}
//     DB.Find(&scenes)
//     for i := 0; i < len(scenes); i++ {
//       scenes[i].discover(locations[i].UserID, lat, lon)
//     }
//   }
// }

func MarkScenesRequest(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}
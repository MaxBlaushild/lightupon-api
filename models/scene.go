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
  PercentDiscovered float64 `sql:"-"`
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
  trips := []Trip{}
  DB.Preload("User").Preload("Scenes.Cards").Order("created_at desc").Where("user_id = ?", userID).Find(&trips)
  for _, trip := range trips {
    for _, scene := range trip.Scenes {
      scene.Trip = trip
      scenes = append(scenes, scene)
    }
  }
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
  isAtNextScene = distanceFromScene < unlockThresholdSmall
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
  // Modifying the radius is necessary because the distanceString below doesn't represent the actual distance in meters, which is more expensive to compute and unnecessary.
  radius = "(" + radius + "^2)*0.000000000080815075"
  distanceString := "((scenes.latitude - " + lat + ")^2.0 + ((scenes.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Where(distanceString + " < " + radius).Order(distanceString + " asc").Limit(numScenes).Find(&scenes)
  for i, _ := range scenes {
    scenes[i].SetPercentDiscovered(userID)
  }
  return
}

func (s *Scene) SetPercentDiscovered(userID uint) (err error) {
  discoveredScene := DiscoveredScene{UserID : userID, SceneID : s.ID}
  err = DB.First(&discoveredScene, discoveredScene).Error; if err == nil {
    s.PercentDiscovered = discoveredScene.PercentDiscovered
  }
  return
}

func MarkScenesRequest(lat string, lon string, userID uint, context string) {
  LogUserLocation(lat, lon, userID, context)
  return
}

func LogUserLocation(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}

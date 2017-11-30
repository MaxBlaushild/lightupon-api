package models

import(
      "strconv"
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/services/imageMagick"
      "io/ioutil"
      "lightupon-api/services/aws"
      "net/http"
      "math"
      "time"
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
  ShareOnFacebook bool
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
  RawScore float64 `sql:"-"`
  TimeVoteScore float64 `sql:"-"`
  SpatialScore float64 `sql:"-"`
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

func GetFollowingScenes(userID uint, page int) (scenes []Scene) {
  limit := 20
  offset := limit * page
  DB.Preload("Trip.User").Preload("Cards").Order("created_at desc").Offset(offset).Limit(limit).Find(&scenes)
  return
}

func GetScenesNearLocation(lat string, lon string, userID uint, radius string, numResults int) (scenes []Scene, err error) {
  // Modifying the radius is necessary because the distanceString below doesn't represent the actual distance in meters, which is more expensive to compute and unnecessary.
  distanceString := "((scenes.latitude - " + lat + ")^2.0 + ((scenes.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Where(distanceString + " < (" + radius + "^2)*0.000000000080815075").Order(distanceString + " asc").Limit(3*numResults).Find(&scenes)

  for i, _ := range scenes {
    scenes[i].SetPercentDiscovered(userID)
  }

  scenes = getTopNScoringScenes(scenes, numResults, userID)

  return
}



func (s *Scene) SetPercentDiscovered(userID uint) (err error) {
  discoveredScene := DiscoveredScene{UserID : userID, SceneID : s.ID}
  err = DB.First(&discoveredScene, discoveredScene).Error; if err == nil {
    s.PercentDiscovered = discoveredScene.PercentDiscovered
  }
  return
}

func LogUserLocation(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}

func (s *Scene) SetTimeVoteScore() {
  timeDiff := time.Now().Sub(s.CreatedAt).Minutes()
  s.TimeVoteScore = s.RawScore / math.Log(timeDiff + 1)
}

func (s *Scene) SetSpatialScore(userLocation UserLocation) {
  distance := CalculateDistance(UserLocation{Latitude: s.Latitude, Longitude: s.Longitude}, userLocation)
  s.SpatialScore = s.TimeVoteScore / math.Log(distance)
}

func getTopNScoringScenes(inputScenes []Scene, n int, userID uint) (scenesToReturn []Scene) {
  location := Location{}
  DB.Where("user_id = ? and context = 'Explore'", userID).Order("created_at desc").First(&location)
  if location.ID == 0 { return; }
  var userLocation UserLocation; userLocation.Latitude = location.Latitude; userLocation.Longitude = location.Longitude // this sucks, we need to refactor

  for i := 0; i < len(inputScenes); i++ {
    inputScenes[i].SetTimeVoteScore()
    inputScenes[i].SetSpatialScore(userLocation)
  }

  var topScoringIndex int
  var topScore float64
  for k := 0; k < n; k++ {
    topScore = 0; topScoringIndex = 0
    for i := 0; i < len(inputScenes); i++ {
      if inputScenes[i].SpatialScore > topScore {
        topScore = inputScenes[i].SpatialScore
        topScoringIndex = i
      }
    }
    scenesToReturn = append(scenesToReturn, inputScenes[topScoringIndex])
    inputScenes = removeSceneFromSlice(inputScenes, topScoringIndex)
  }
  return
}

func removeSceneFromSlice(inputScenes []Scene, indexToRemove int) (scenesToReturn []Scene) {
    for i := 0; i < len(inputScenes); i++ {
        if i != indexToRemove {
            scenesToReturn = append(scenesToReturn, inputScenes[i])
        }
    }
    return
}
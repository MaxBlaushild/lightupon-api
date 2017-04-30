package models

import(
// "github.com/davecgh/go-spew/spew"
      "strconv"
      "fmt"
      "github.com/jinzhu/gorm"
      "github.com/nfnt/resize"
      "lightupon-api/services/aws"
      "net/http"
      "image"
      "image/jpeg"
      "log"
      "bytes"
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

type ExposedScene struct {
  gorm.Model
  UserID uint
  SceneID uint
  Blur float64
  Unlocked bool
}

const unlockThresholdSmall = 0.06
const unlockThresholdLarge = 0.4

func (s *Scene) AfterCreate(tx *gorm.DB) (err error) {
  err = s.uploadAndSetPins()
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

func (s *Scene) downloadAndCropImage() (pinImageBinary []byte, selectedPinImageBinary []byte){
  response, _ := http.Get(s.BackgroundUrl)
  image, _, _ := image.Decode(response.Body)
  defer response.Body.Close()

  pinImage := resize.Resize(40, 0, image, resize.Lanczos3)
  pinBuffer := new(bytes.Buffer)
  if err := jpeg.Encode(pinBuffer, pinImage, nil); err != nil {
    log.Println("unable to encode image.")
  }
  pinImageBinary = pinBuffer.Bytes()

  selectedPinImage := resize.Resize(80, 0, image, resize.Lanczos3)
  selectedPinBuffer := new(bytes.Buffer)
  if err := jpeg.Encode(selectedPinBuffer, selectedPinImage, nil); err != nil {
    log.Println("unable to encode image.")
  }
  selectedPinImageBinary = selectedPinBuffer.Bytes()

  return
}

func (s *Scene) getAssetName(name string) string {
  return "/scenes/" + fmt.Sprint(s.ID) + "/" + name
}

func (s *Scene) uploadAndSetPins() (err error) {
  pinBinary, selectedPinBinary := s.downloadAndCropImage()
  s.PinUrl, err = s.uploadPin(pinBinary, "pin")
  s.SelectedPinUrl, err = s.uploadPin(selectedPinBinary, "selectedPin")
  return
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

func GetScenesNearLocation(lat string, lon string, userID uint) (scenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Order("((scenes.latitude - " + lat + ")^2.0 + ((scenes.longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc").Limit(50).Find(&scenes)

  for i := 0; i < len(scenes); i++ {
    scenes[i].discover(userID, lat, lon)
  }

  // possiblyRecomputeAllDiscovery(lat, lon, userID)

  return
}

// There's not a lot of awesome abstraction here, but this could get computationally expensive so I'm trying to optimize for speed. It also requires some splainin so read the comments.
func (scene *Scene) discover(userID uint, userLat string, userLon string) {
  // Try to get the current blur level
  oldExposedScene := ExposedScene{UserID : userID, SceneID : scene.ID}
  DB.First(&oldExposedScene, oldExposedScene)

  if (oldExposedScene.Unlocked) { // UU, UB, UL
    // If we found a record and it's fully unlocked, then we're done. Just set the properties - no need to persist anything.
    scene.Blur = 0.0
    scene.Hidden = false
  } else { // BU, BB, BL, LU, LB, LL
    userLatFloat, _ := strconv.ParseFloat(userLat, 64)
    userLonFloat, _ := strconv.ParseFloat(userLon, 64)
    distanceFromScene := CalculateDistance(UserLocation{Latitude: userLatFloat, Longitude: userLonFloat}, UserLocation{Latitude: scene.Latitude, Longitude: scene.Longitude})
    newBlur := calculateBlur(distanceFromScene)
    if (distanceFromScene < unlockThresholdSmall) { // LU, BU
      // Moving to Unlocked from unlockThresholdLarge non-Unlocked state, so we need to return and persist the new shit
      scene.Blur = 0.0
      scene.Hidden = false
      oldExposedScene.upsertExposedScene(newBlur, scene.ID, userID, false)
    } else {
      if (newBlur > oldExposedScene.Blur) { // MM(change), LM
        // save the new blur and return that new shit
        scene.Blur = newBlur
        scene.Hidden = true
        oldExposedScene.upsertExposedScene(newBlur, scene.ID, userID, true)
      } else { // MM(static), ML, LL
        // save nothing and return the old shit
        scene.Blur = oldExposedScene.Blur
        scene.Hidden = true
      }
    }
  }
}

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

func (exposedScene *ExposedScene) upsertExposedScene(newBlur float64, sceneID uint, userID uint, hidden bool) {
  if (exposedScene.ID == 0) {
    DB.Create(&ExposedScene{UserID : userID, SceneID : sceneID, Blur : newBlur, Unlocked : !hidden})
  } else {
    DB.Model(&exposedScene).Update("Blur", newBlur).Update("Unlocked", !hidden)
  }
}

// Under construction
// func recomputeAllDiscovery() {
//   exposedScene := ExposedScene{}

//   locations := []Location{}
//   DB.Find(&locations)
//   for i := 0; i < len(locations); i++ {
//     lat := strconv.FormatFloat(locations[i].Latitude, 'E', -1, 64)
//     lon := strconv.FormatFloat(locations[i].Longitude, 'E', -1, 64)
//     // Here we should grab ALL scenes out of the database and iterate over them
//     scenes := []Scene{}
//     DB.Find(&scenes)
//     for i := 0; i < len(scenes); i++ {
//       scene.discover(userID uint, userLat string, userLon string)
      
//     }

//     actuallyUpdateUserDarknessState(lat, lon, locations[i].UserID)
//   }
// }

func MarkScenesRequest(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}
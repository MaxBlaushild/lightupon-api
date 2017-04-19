package models

import(
      "strconv"
      "fmt"
      "github.com/jinzhu/gorm"
      "lightupon-api/services/aws"
      "lightupon-api/services/imageMagick"
      "io/ioutil"
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
}

type ExposedScene struct {
  gorm.Model
  UserID uint
  SceneID uint
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

func IndexScenes() (scenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Order("created_at desc").Find(&scenes)
  for i := 0; i < len(scenes); i++ {
    scenes[i].darken()
  }
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

func (s *Scene) getAssetName(name string) string {
  return "/scenes/" + fmt.Sprint(s.ID) + "/" + name
}

func (s *Scene) SetPins() (err error) {
  imageBinary := s.DownloadImage()
  pinBinary := imageMagick.CropPin(imageBinary, "40x40!")
  selectedPinBinary := imageMagick.CropPin(imageBinary, "80x80!")
  s.PinUrl, err = s.uploadPin(pinBinary, "pin")
  s.SelectedPinUrl, err = s.uploadPin(selectedPinBinary, "selectedPin")
  fmt.Println(err)
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

  userNeighborhood := getNeighborhoodIDForLocation(lat, lon)

  for i := 0; i < len(scenes); i++ {
    sceneNeighborhood := getNeighborhoodIDForLocation(strconv.FormatFloat(scenes[i].Latitude, 'f', 6, 64), strconv.FormatFloat(scenes[i].Longitude, 'f', 6, 64))
    if (scenes[i].hiddenFromUser(userID) && (sceneNeighborhood != userNeighborhood)) {
      scenes[i].darken()
    }
  }

  scenes = append(scenes, getNeighborhoodScenes()...)

  return
}

func (scene *Scene) darken() {
  scene.BackgroundUrl = "http://www.solidbackgrounds.com/images/2560x1440/2560x1440-black-solid-color-background.jpg"
  scene.Name = ""
  if (len(scene.Cards) > 0) {
    scene.Cards[0].Caption = ""
    // scene.Cards[0].ImageURL = "http://www.solidbackgrounds.com/images/2560x1440/2560x1440-black-solid-color-background.jpg" // Let's not darken this one just yet
  }
  return
}

func getNeighborhoodScenes() (neighborhoodScenes []Scene) {
  DB.Preload("Trip.User").Preload("Cards").Preload("SceneLikes").Where("scenes.Name IN ('Brookline', 'Fenway', 'Back Bay', 'South End', 'Seaport', 'Downtown', 'Cambridge')").Find(&neighborhoodScenes)
  return
}

func (s *Scene) hiddenFromUser(userID uint) bool {
  exposedScenes := []ExposedScene{}
  sql := `SELECT * FROM exposed_scenes
          WHERE user_id = ` + strconv.Itoa(int(userID)) + `
          AND scene_id = ` + strconv.Itoa(int(s.ID)) + `;`
  DB.Raw(sql).Scan(&exposedScenes)

  return len(exposedScenes) == 0
}

func GetScenesVeryNearLocation(lat string, lon string) (scenes []Scene) {
  sql := `SELECT * FROM scenes
          WHERE ((latitude - ` + lat + `)^2.0 + ((longitude - ` + lon + `)* cos(latitude / 57.3))^2.0) < 0.000003;`; // 0.000003 is about one block on newbury st
  DB.Raw(sql).Scan(&scenes)
  return
}

func MarkScenesRequest(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}
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
  // TODO: Stick all of our data access in one file so Max doesn't ever have to look at it
  sql := `SELECT * FROM scenes
          INNER JOIN follows ON scenes.user_id = follows.followed_user_id
          WHERE following_user_id = ` + strconv.Itoa(int(userID)) + `
          AND ((latitude - ` + lat + `)^2.0 + ((longitude - ` + lon + `)* cos(latitude / 57.3))^2.0) < 1
          ORDER BY ((latitude - ` + lat + `)^2.0 + ((longitude - ` + lon + `)* cos(latitude / 57.3))^2.0) asc;`

  DB.Raw(sql).Scan(&scenes)

  return
}

func GetScenesNearLocation(lat string, lon string, userID uint) (scenes []Scene) {
  // TODO: Stick all of our data access in one file so Max doesn't ever have to look at it
  sql := `SELECT * FROM scenes
          ORDER BY ((latitude - ` + lat + `)^2.0 + ((longitude - ` + lon + `)* cos(latitude / 57.3))^2.0) asc
          LIMIT 20;`

  DB.Raw(sql).Scan(&scenes)

  for i := 0; i < len(scenes); i++ {
    fmt.Println(scenes[i].hiddenFromUser(userID))
    if scenes[i].hiddenFromUser(userID) {
      fmt.Println("do cool stuff")
      scenes[i].BackgroundUrl = "http://wallpaperrs.com/uploads/nature/earth-moon-night-field-stupendous-wallpaper-88356-142977423728.jpg"
      scenes[i].Name = "Darklands..."
    }
  }


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


func MarkScenesRequest(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := Location{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}
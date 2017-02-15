package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"os"
	"time"
  "net/http"
  "strconv"
  "lightupon-api/services/redis"
  "lightupon-api/services/googleMaps"
)

type User struct {
	gorm.Model
	FacebookId string `gorm:"type:varchar(100);unique_index"`
	Email string `gorm:"type:varchar(100);unique_index"`
	FirstName string
	ProfilePictureURL string
	FullName string
	DeviceID string
  Devices []Device
	Token string
  Scenes []Scene
	Parties []Party `gorm:"many2many:partyusers;"`
	Lit bool
	Trips []Trip
  SceneLikes []SceneLike
	Location UserLocation `gorm:"-"`
	Follows []Follow `gorm:"ForeignKey:FollowingUserID"`
  NumberOfFollowers int `sql:"-"`
  NumberOfTrips int `sql:"-"`
  Following bool `sql:"-"`

}

const threshold float64 = 0.05 // 0.05 km = 50 meters

func (u *User) BeforeCreate() (err error) {
  u.Token = createToken(u.FacebookId)
  return
}

func (u *User) AddTrip(trip *Trip) (err error) {
  err = DB.Model(&u).Association("Trips").Append(&trip).Error
  err = DB.Model(&u).Association("Scenes").Append(&trip.Scenes[0]).Error
  return
}

func (u *User) IsFollowing(user *User) bool {
  var count int
  DB.Model(&Follow{}).Where("followed_user_id = ? AND following_user_id = ?", user.ID, u.ID).Count(&count)
  return (count > 0)
}

func (u *User) PopulateIsFollowing(user *User) {
  u.Following = user.IsFollowing(u)
}

func GetUserByID(userID string) (user User){
  DB.Where("id = ?", userID).First(&user)
  user.PopulateNumberOfFollowers()
  user.PopulatingNumberOfTrips()
  return
}

func (u *User) GetNumberOfTrips() int {
  count := DB.Model(&u).Association("Trips").Count()
  return count
}

func (u *User) PopulatingNumberOfTrips() {
  u.NumberOfTrips = u.GetNumberOfTrips()
}

func (u *User) PopulateNumberOfFollowers() {
  u.NumberOfFollowers = u.GetFollowerCount()
}

func (u *User) GetFollowerCount() (count int) {
  DB.Model(&Follow{}).Where("followed_user_id = ?", u.ID).Count(&count)
  return
}

func UpsertUser(user *User) {
	DB.Where("facebook_id = ?", user.FacebookId).Assign(user).FirstOrCreate(&user)

	if !DB.NewRecord(user) {
		DB.Save(user)
	}
}

func (u *User) SetUserLikenessOfScenes(scenes []Scene) {
  for i, scene := range scenes {
    scenes[i].Liked = scene.UserHasLiked(u)
  }
}

func FindUsers(query string) (users []User) {
	fuzzyQuery := "%" + query
	fuzzyQuery += "%"
	DB.Where("full_name ILIKE ?", fuzzyQuery).Find(&users)
	return
}

func createToken(facebookId string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["facebookId"] = facebookId
	token.Claims["exp"] = time.Now().Add(time.Hour * 72000).Unix() // For now, set tokens to expire never
	signingSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(signingSecret); if err != nil {
  	fmt.Println(err)
  }
  return tokenString
}

func (u *User) ActiveParty() (activeParty Party) {
  DB.Model(&u).Preload("Trip.Scenes.Cards").Where("active = true").Association("Parties").Find(&activeParty)
  return
}

func RefreshTokenByFacebookId(facebookId string) string {
	user := User{}
	token := createToken(facebookId)
	DB.Model(&user).Where("facebookId = ?", facebookId).Update("token", token)
  return token
}

func (u *User) IsAtScene(scene Scene)(isAtNextScene bool) {
	sceneLocation := UserLocation{Latitude:scene.Latitude, Longitude: scene.Longitude}
	distanceFromScene := CalculateDistance(sceneLocation, u.Location)
	isAtNextScene = distanceFromScene < threshold
	return
}

func (u *User) AddLocationToCurrentTrip(location Location)(err error) {
  trip := Trip{}
  DB.Where("user_id = ? AND active = true", u.ID).First(&trip)

  if (trip.ID > 0) {
    err = DB.Model(&trip).Association("Locations").Append(location).Error
    SaveCurrentLocationToRedis(u.FacebookId, location)
  }

  return
}

func (u *User) ActiveTrip()(trip Trip) {
  DB.Preload("Scenes.Cards").Where("user_id = ? AND active = true", u.ID).First(&trip)
  return
}

func (u *User) SetUserLocationFromRequest(r *http.Request) {
  query := r.URL.Query()

  location := UserLocation{}
  lat, _ := strconv.ParseFloat(query["lat"][0], 64)
  lon, _ := strconv.ParseFloat(query["lon"][0], 64)
  location.Latitude = lat
  location.Longitude = lon

  u.Location = location
}

func (u *User) GetActiveSceneOrSuggestion() (scene Scene) {
  activeTrip := u.ActiveTrip()

  lengthOfScenes := len(activeTrip.Scenes)

  if (lengthOfScenes == 0) {
    return u.GetSuggestedScene()
  }

  activeScene := activeTrip.Scenes[lengthOfScenes - 1]

  if (u.IsAtScene(activeScene)) {
    return activeScene
  } else {
    return u.GetSuggestedScene()
  }

  return
}

func (u *User) GetSuggestedScene() (scene Scene) {
  place := googleMaps.GetPrettyPlace(u.Location.Latitude, u.Location.Longitude)
  scene.FormattedAddress = place["FormattedAddress"]
  scene.StreetNumber = place["street_number"]
  scene.Route = place["route"]
  scene.Neighborhood = place["neighborhood"]
  scene.Locality = place["locality"]
  scene.AdministrativeLevelTwo = place["administrative_area_level_2"]
  scene.AdministrativeLevelOne = place["administrative_area_level_1"]
  scene.Country = place["country"]
  scene.PostalCode = place["postal_code"]
  scene.GooglePlaceID = place["PlaceID"]
  return scene
}

func (u *User) DeactivateTrips() {
	DB.Model(&Trip{}).Where("active = true AND user_id = ?", u.ID).Update("active", false)
}

func (u *User) Light(location Location)(err error) {
	tx := DB.Begin()

  if err := tx.Model(&u).Update("lit", true).Error; err != nil {
    tx.Rollback()
    return err
  }

  trip := Trip{
                Active: true,
  							ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/e/e4/Stourhead_garden.jpg",
  							Description: "This is the song that never ends.",
  							Details: "And it goes on and on my friends.",
  						}

  if err := tx.Model(&u).Association("Trips").Append(&trip).Error; err != nil {
  	fmt.Println(err)
     tx.Rollback()
     return err
  }

  tx.Commit()
  return nil
}

func (u *User) Extinguish(location Location)(err error) {
	DB.Model(&u).Update("lit", false)
  u.DeactivateTrips()
  redis.DeleteRedisKey("currentLocation_" + u.FacebookId)
  return nil
}

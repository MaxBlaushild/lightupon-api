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
  "lightupon-api/live"
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

func (user *User) Explore() (err error)  {
  scenes, err := GetScenesNearLocation(fmt.Sprintf("%.6f", user.Location.Latitude), fmt.Sprintf("%.6f", user.Location.Longitude), user.ID, "0.01", 100)
  for i := 0; i < len(scenes); i++ {
    user.Discover(scenes[i])
  }
  return
}


// There's not a lot of awesome abstraction here, but this could get computationally expensive so I'm trying to optimize for speed. It also requires some splainin so read the comments.
func (user *User) Discover(scene Scene) {
  // Try to get the current blur level
  oldExposedScene := ExposedScene{UserID : user.ID, SceneID : scene.ID}
  DB.First(&oldExposedScene, oldExposedScene)

  if (oldExposedScene.Unlocked) { // UU, UB, UL
    // If we found a record and it's fully unlocked, then we're done. Just set the properties - no need to persist anything.
    scene.Blur = 0.0
    scene.Hidden = false
  } else { // BU, BB, BL, LU, LB, LL
    distanceFromScene := CalculateDistance(user.Location, UserLocation{Latitude: scene.Latitude, Longitude: scene.Longitude})
    newBlur := calculateBlur(distanceFromScene)
    if (distanceFromScene < unlockThresholdSmall) { // LU, BU
      // Moving to Unlocked from unlockThresholdLarge non-Unlocked state, so we need to return and persist the new shit
      scene.Blur = 0.0
      scene.Hidden = false
      oldExposedScene.upsertExposedScene(newBlur, scene.ID, user.ID, false)
    } else {
      if (newBlur < oldExposedScene.Blur) { // MM(change), LM
        // save the new blur and return that new shit
        scene.Blur = newBlur
        scene.Hidden = true
        oldExposedScene.upsertExposedScene(newBlur, scene.ID, user.ID, true)
      } else { // MM(static), ML, LL
        // save nothing and return the old shit
        scene.Blur = oldExposedScene.Blur
        scene.Hidden = true
        fmt.Println("i dont do anything")
      }
    }
  }
}

func (u *User) StartParty(tripID uint) (party Party, err error) {
  party = Party{ TripID: tripID }
  err = DB.Model(&u).Association("Parties").Append(&party).Error
  party.LoadTrip()
  live.Hub.AddUserToParty(u.ID, party.Passcode)
  return
}

func (u *User) AddTrip(trip *Trip) (err error) {
  err = DB.Model(&Trip{}).Where("user_id = ?", u.ID).Update("Active", false).Error
  err = DB.Model(&u).Association("Trips").Append(trip).Error
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
  activeParty.SyncWithLive()
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

func (u *User) Light() (err error) {
  err = DB.Model(&u).Update("lit", true).Error
  return
}

func (u *User) Extinguish()(err error) {
	DB.Model(&u).Update("lit", false)
  u.DeactivateTrips()
  redis.DeleteRedisKey("currentLocation_" + u.FacebookId)
  return nil
}



type UserStats struct {
  Name string
  NumExposedScenes int
  NumSceneGets int
  NumUpvotesGiven int
}

func GetUserStats() (stats []UserStats) {
  sql := `SELECT u.first_name AS name,
                 es.num AS num_exposed_scenes,
                 l.num AS num_scene_gets,
                 sl.num AS num_upvotes_given
          FROM users u
          LEFT OUTER JOIN (SELECT user_id, COUNT(*) AS num FROM scene_likes GROUP BY user_id) sl
          ON u.id = sl.user_id
          LEFT OUTER JOIN (SELECT user_id, COUNT(*) AS num FROM exposed_scenes GROUP BY user_id) es
          ON u.id = es.user_id
          LEFT OUTER JOIN (SELECT user_id, COUNT(*) AS num FROM locations GROUP BY user_id) l
          ON u.id = l.user_id;`;
  DB.Raw(sql).Scan(&stats)
  return
}


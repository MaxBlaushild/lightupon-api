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
  "lightupon-api/services/facebook"
  "lightupon-api/services/twitter"
  // "github.com/kr/pretty"
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
  Posts []Post
  FacebookToken string
  TwitterKey string
  TwitterSecret string
	Lit bool
	Trips []Trip
	Location UserLocation `gorm:"-"`
  ActualLocation Location
	Follows []Follow `gorm:"ForeignKey:FollowingUserID"`
  NumberOfFollowers int `sql:"-"`
  NumberOfTrips int `sql:"-"`
  Following bool `sql:"-"`
  ManaTotal int `sql:"-"`
}

func (u *User) BeforeCreate() (err error) {
  u.Token = createToken(u.FacebookId)
  return
}

func (u *User) PostToFacebook(p *Post) (err error) {
  fbUser := facebook.User{
    ID: u.FacebookId,
    AccessToken: u.FacebookToken,
  }

  post := facebook.Post {
    Message: p.Caption,
    PictureUrl: p.ImageUrl,
    Link: p.ImageUrl,
  }

  err = facebook.CreatePost(fbUser, post)
  return
}

func (u *User) PostToTwitter(p *Post) (err error) {
  twitterUser := twitter.User{
    AccessToken: u.TwitterKey,
    AccessTokenSecret: u.TwitterSecret,
  }

  postImageBinary, err := DownloadImage(p.ImageUrl)
  media, err := twitter.PostMedia(twitterUser, postImageBinary); if err != nil {
    return
  }

  status := twitter.Status{
    Status: p.Caption,
    Lat: p.Latitude,
    Long: p.Longitude,
    MediaID: media.MediaIDString,
  }

  err = twitter.PostStatus(twitterUser, status)
  return
}

func (user *User) Update(updates User) (err error) {
  err = DB.Model(&user).Update(updates).Error
  return
}

func (user *User) Explore() (err error)  {
  latString := fmt.Sprintf("%.6f", user.Location.Latitude)
  lonString := fmt.Sprintf("%.6f", user.Location.Longitude)
  LogUserLocation(latString, latString, user.ID, "Explore")
  posts, err := GetPostsNearLocation(latString, lonString, user.ID, fmt.Sprintf("%.6f", unlockThresholdLarge), 100)
  for i := 0; i < len(posts); i++ {
    user.Discover(&posts[i])
  }
  return
}

func (user *User) Discover(post *Post) {
  currentDiscoveredPost := GetCurrentDiscoveredPost(user.ID, post.ID)
  if currentDiscoveredPost.NotFullyDiscovered() {
    currentDiscoveredPost.UpdatePercentDiscovered(user, post)
  }
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

func UpsertUser(userToUpsert *User) {
  user := User{}
	DB.Where(User{ FacebookId: userToUpsert.FacebookId}).FirstOrCreate(&user)
  userToUpsert.ID = user.ID
	DB.Save(userToUpsert)
}

func UserIsBlackListed(token string) bool {
  blacklistUser := BlacklistUser{Token : token}
  DB.First(&blacklistUser)
  return blacklistUser.ID != 0
}

func (u *User) SetUserLikenessOfScenes(scenes []Scene) {
  // for i, scene := range scenes {
  //   scenes[i].Liked = scene.UserHasLiked(u)
  // }
}

func FindUsers(query string) (users []User) {
	fuzzyQuery := "%" + query
	fuzzyQuery += "%"
	DB.Where("full_name ILIKE ?", fuzzyQuery).Find(&users)
	return
}

func createToken(facebookId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "facebookId": facebookId,
    "exp": time.Now().Add(time.Hour * 72000).Unix(),
  })
	signingSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(signingSecret); if err != nil {
  	fmt.Println(err)
  }
  return tokenString
}

func RefreshTokenByFacebookId(facebookId string) string {
	user := User{}
	token := createToken(facebookId)
	DB.Model(&user).Where(User{ FacebookId: facebookId}).Update("token", token)
  return token
}

func (u *User) IsAtScene(scene Scene)(isAtNextScene bool) {
	sceneLocation := UserLocation{Latitude:scene.Latitude, Longitude: scene.Longitude}
	distanceFromScene := CalculateDistance(sceneLocation, u.Location)
	isAtNextScene = distanceFromScene < unlockThresholdSmall
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
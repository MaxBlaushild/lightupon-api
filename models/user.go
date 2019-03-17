package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"os"
	"time"
  "net/http"
  "strconv"
  "lightupon-api/services/facebook"
  "lightupon-api/services/twitter"
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
  Posts []Post
  FacebookToken string
  TwitterKey string
  TwitterSecret string
	Location Location `gorm:"-"`
  ActualLocation Location
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

func (u *User) TrackedQuests () (quests []Quest, err error) {
  trackedQuests := []TrackedQuest{}

  // surprised this works, but ok!!!!!
  err = DB.Preload("Quest.Posts").Preload("Quest.QuestProgress").Where("user_id = ?", u.ID).Find(&trackedQuests).Error; if err !=  nil {
    return
  }

  for _, trackedQuest := range trackedQuests {
    quest := trackedQuest.Quest
    quests = append(quests, quest)
  }

  return
}

func (u *User) TrackQuest(questID uint) (err error) {
  trackedQuest := TrackedQuest{}
  err = DB.FirstOrCreate(&trackedQuest, TrackedQuest{ QuestID: questID, UserID: u.ID }).Error
  return
}

func (u *User) UntrackQuest(questID uint) (err error) {
  trackedQuest := TrackedQuest{}
  err = DB.Where("user_id = ? AND quest_id = ?", u.ID, questID).Delete(&trackedQuest).Error
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

func GetUserByID(userID string) (user User){
  DB.Where("id = ?", userID).First(&user)
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

func (u *User) SetLocationFromRequest(r *http.Request) {
  query := r.URL.Query()

  location := Location{}
  lat, _ := strconv.ParseFloat(query["lat"][0], 64)
  lon, _ := strconv.ParseFloat(query["lon"][0], 64)
  location.Latitude = lat
  location.Longitude = lon

  u.Location = location
}
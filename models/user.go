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

func (u *User) ActiveQuests() (quests []Quest, err error) {
  completedPosts := []DiscoveredPost{}

  err = DB.Where(DiscoveredPost{ 
    UserID: u.ID,
    Completed: true,
  }).Preload("Post").Find(&completedPosts).Error; if err != nil {
    return
  }

  postCounts := map[uint]uint{}
  uniqueIDs := []uint{}

  for _, discoveredPost := range completedPosts {
    questID := discoveredPost.Post.QuestID

    _, isIncluded := postCounts[questID]; if isIncluded {
      postCounts[questID] += 1
    } else {
      postCounts[questID] = 1
      uniqueIDs = append(uniqueIDs, questID)
    }
  }

  for _, id := range uniqueIDs {
    quest := Quest{}
    err = DB.Preload("Posts").First(&quest, id).Error; if err != nil { return }
    quest.QuestProgress.CompletedPosts = postCounts[id]
    quests = append(quests, quest)
  }

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
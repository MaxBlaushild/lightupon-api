package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

type User struct {
	gorm.Model
	FacebookId string `gorm:"type:varchar(100);unique_index"`
	Email string `gorm:"type:varchar(100);unique_index"`
	FirstName string
	ProfilePictureURL string
	FullName string
	DeviceID string
	Token string
	Parties []Party `gorm:"many2many:partyusers;"`
	Lit bool
	Trips []Trip
	Location UserLocation `gorm:"-"`
}

const threshold float64 = 0.05 // 0.05 km = 50 meters

func (u *User) BeforeCreate() (err error) {
  u.Token = createToken(u.FacebookId)
  return
}

func UpsertUser(user User) {
	DB.Where("facebook_id = ?", user.FacebookId).Assign(user).FirstOrCreate(&user)
	
	if !DB.NewRecord(user) {
		DB.Save(user)
	}
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
	DB.Where("user_id = ?", u.ID).Last(&trip)
	err = DB.Model(&trip).Association("Locations").Append(location).Error
	return
}

func (u *User) Light()(err error) {
	tx := DB.Begin()

  if err := tx.Model(&u).Update("lit", true).Error; err != nil {
    tx.Rollback()
    return err
  }

  trip := Trip{ Title: "LOG DATE: TANGO",
  							ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/e/e4/Stourhead_garden.jpg",
  							Description: "This is the song that never ends.",
  							Details: "And it goes on and on my friends.",
  						}

  if err := tx.Model(&u).Association("trips").Append(trip).Error; err != nil {
     tx.Rollback()
     return err
  }

  tx.Commit()
  return nil
}

func (u *User) Extinguish()(err error) {
	err = DB.Model(&u).Update("lit", false).Error
	return
}

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
	FacebookId string `gorm:"index"`
	Email string
	DeviceID string
	Token string
	Parties []Party `gorm:"many2many:partyusers;"`
	Location Location `gorm:"-"`
}

const threshold float64 = 0.05 // 0.05 km = 50 meters

func (u *User) BeforeCreate() (err error) {
  u.Token = createToken(u.FacebookId)
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

func (u *User) IsAtScene(scene Scene)(isAtNextScene bool) {
	sceneLocation := Location{scene.Latitude, scene.Longitude}
	distanceFromScene := CalculateDistance(sceneLocation, u.Location)
	isAtNextScene = distanceFromScene < threshold
	return
}

func RefreshTokenByFacebookId(facebookId string) string {
	user := User{}
	token := createToken(facebookId)
	DB.Model(&user).Where("facebookId = ?", facebookId).Update("token", token)
  return token
}

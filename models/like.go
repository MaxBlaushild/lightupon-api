package models

import(
      "github.com/jinzhu/gorm"
      // "github.com/davecgh/go-spew/spew"
      //  "fmt"
      )

type Like struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  BookmarkID uint `gorm:"not null"`
}

func GetLikesForUser(userID uint) (likes []Like){
  DB.Where("user_id = $1", userID).Find(&likes)
  return
}
package models

import(
      "github.com/jinzhu/gorm"
      // "github.com/davecgh/go-spew/spew"
      //  "fmt"
      )

type FuckThis struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  BookmarkID uint `gorm:"not null"`
}

func GetFuckThisesForUser(userID uint) (fuckThises []FuckThis){
  DB.Where("user_id = $1", userID).Find(&fuckThises)
  return
}
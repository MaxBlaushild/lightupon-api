package models

import(
      "github.com/jinzhu/gorm"
      )


type Device struct {
    gorm.Model
    UserID uint
    User User
    DeviceToken string
}
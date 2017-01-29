package models

import(
      "github.com/jinzhu/gorm"
      )


type DeviceToken struct {
    gorm.Model
    UserID uint
    DeviceToken string
}
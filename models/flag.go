package models

import(
      "github.com/jinzhu/gorm"
      )

type Flag struct {
  gorm.Model
  UserID uint
  PostID uint
  Description string
}
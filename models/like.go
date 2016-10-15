package models

import(
      "github.com/jinzhu/gorm"
      )

type Like struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  BookmarkID uint `gorm:"not null"`
}
package models

import(
      "github.com/jinzhu/gorm"
      )

type Bookmark struct {
  gorm.Model
  Title string `gorm:"not null"`
  URL string
}
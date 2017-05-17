package models

import(
      "github.com/jinzhu/gorm"
      )

type BlacklistUser struct {
  gorm.Model
  Token string
}
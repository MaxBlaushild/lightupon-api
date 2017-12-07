package models

import(
      "github.com/jinzhu/gorm"
      )

type Entity struct {
  gorm.Model
  UserID uint
  SceneID uint
  Description string
}
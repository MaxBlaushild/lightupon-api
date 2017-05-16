package models

import(
      "github.com/jinzhu/gorm"
      )

type Flag struct {
  gorm.Model
  UserID uint
  SceneID uint
}
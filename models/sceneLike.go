package models

import(
      "github.com/jinzhu/gorm"
      // "github.com/davecgh/go-spew/spew"
      // "fmt"
      )

type SceneLike struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  User User
  Scene Scene
  SceneID uint
}
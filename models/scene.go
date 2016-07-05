package models

import(
      "github.com/jinzhu/gorm"
      )

type Scene struct {
  gorm.Model
  Name string
  Latitude float64
  Longitude float64
  TripID uint
  BackgroundUrl string
  SceneOrder uint
  Cards []Card
}

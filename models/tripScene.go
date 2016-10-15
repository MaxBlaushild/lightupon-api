package models

import(
      "github.com/jinzhu/gorm"
      )

type TripScene struct {
  gorm.Model
  Title string
  Description string
  ImageUrl string
  Distance float32
  Latitude float64
  Longitude float64
  EstimatedTime int
  Owner int
  Scenes []Scene
}
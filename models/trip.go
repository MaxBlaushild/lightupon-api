package models

import(
      "github.com/jinzhu/gorm"
      )

type Trip struct {
  gorm.Model
  Title string
  Description string
  ImageUrl string
  Distance float32
  Latitude float64
 	Longitude float64
 	EstimatedTime int
}


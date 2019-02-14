package models

import(
      "github.com/jinzhu/gorm"
      )

type Trip struct {
  gorm.Model
  Title string `gorm:"not null"`
  Description string `gorm:"not null"`
  Details string
  ImageUrl string `gorm:"not null"`
  Distance float32
  EstimatedTime int
  UserID uint
  User User
  Scenes []Scene
  Locations []Location
  Active bool
  Constellation []ConstellationPoint
  UserHasLikedTrip bool `sql:"-"`
  TotalLikes int `sql:"-"`
}

type ConstellationPoint struct {
  DeltaY float64
  DistanceToPreviousPoint float64
}
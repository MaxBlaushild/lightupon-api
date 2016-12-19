package models

import("github.com/jinzhu/gorm")

type Comment struct {
  gorm.Model
  User User
  UserID uint
  Text string
  TripID uint
  Trip Trip
  SceneID uint
  Scene Scene
}
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
  CardID uint
  Card Card
}

func GetCommentsForScene(sceneID int)(comments []Comment) {
	DB.Where("scene_id = ?", sceneID).Preload("User").Find(&comments)
	return
}

func GetCommentsForTrip(tripID int)(comments []Comment) {
	DB.Where("trip_id = ?", tripID).Preload("User").Find(&comments)
	return
}

func GetCommentsForCard(cardID int)(comments []Comment) {
  DB.Where("card_id = ?", cardID).Preload("User").Find(&comments)
  return
}

func (trip *Trip) LoadCommentsForTrip() {
  trip.Comments = GetCommentsForTrip(int(trip.ID))
  return
}
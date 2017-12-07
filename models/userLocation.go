package models

import("strconv")

type UserLocation struct {
	Latitude float64
	Longitude float64
  UserID uint
	Context string
	Course float64
  Accuracy float64
  Floor int
}

func LogUserLocation(lat string, lon string, userID uint, context string) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)
  location := UserLocation{UserID:userID, Latitude: latFloat, Longitude: lonFloat, Context: context}
  DB.Create(&location)
  return
}
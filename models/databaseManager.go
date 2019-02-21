package models

import (
	      "github.com/jinzhu/gorm"
)


type DatabaseAccessor interface {
  GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error)
}

type databaseManager struct {
  DB *gorm.DB
}

func CreateNewDatabaseManager (DB *gorm.DB) (databaseManager databaseManager) {
  databaseManager.DB = DB
  return
}

func (databaseManager databaseManager) GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error) {
  distanceString := "((posts.latitude - " + lat + ")^2.0 + ((posts.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  whereClause := distanceString + " < (" + radius + "^2)*0.000000000080815075"
  whereClause += " AND QuestOrder = 1"
  orderClause := distanceString + " asc"
  DB.Preload("Pin").Preload("User").Where(whereClause).Order(orderClause).Limit(numResults).Find(&posts)
  return
}

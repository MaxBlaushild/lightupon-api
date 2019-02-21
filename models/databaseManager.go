package models

import (
	      "github.com/jinzhu/gorm"
        // "github.com/davecgh/go-spew/spew"
)


type DatabaseAccessor interface {
  GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error)
  GetQuestOrderForLastCompletedPostInEachQuest(userID uint) (results []struct{QuestID uint; MaxQuestOrder uint;}, err error)
  FindNearbyPostInQuestWithParticularQuestOrder(lat string, lon string, radius string, questID uint, questOrder uint) (post Post, err error)
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
  databaseManager.DB.Preload("Pin").Preload("User").Where(whereClause).Order(orderClause).Limit(numResults).Find(&posts)
  return
}

func (databaseManager databaseManager) GetQuestOrderForLastCompletedPostInEachQuest(userID uint) (results []struct{QuestID uint; MaxQuestOrder uint;}, err error) {
  query := `SELECT p.quest_id, MAX(quest_order) AS max_quest_order
            FROM discovered_posts dp 
            JOIN posts p ON 
              dp.user_id = ? AND 
              dp.post_id = p.id AND
              dp.Completed = true
            GROUP BY p.quest_id`

  databaseManager.DB.Raw(query, userID).Scan(&results)

  return
}

func (databaseManager databaseManager) FindNearbyPostInQuestWithParticularQuestOrder(lat string, lon string, radius string, questID uint, questOrder uint) (post Post, err error) {
  whereClause := `((posts.latitude - ?)^2.0 + ((posts.longitude - ?)* cos(latitude / 57.3))^2.0) < (?^2)*0.000000000080815075
                      AND quest_id = ?
                      AND quest_order = ?`
  databaseManager.DB.Where(whereClause, lat, lon, radius, questID, questOrder).First(&post)

  // databaseManager.DB.Where("WHERE quest_id = 1").First(&post)
  // databaseManager.DB.First(&post)
  return
}
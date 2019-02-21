package models

import (
  "github.com/jinzhu/gorm"
)

type DatabaseAccessor interface {
  GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error)
  GetQuestOrderForLastCompletedPostInEachQuest(userID uint) (results []struct{QuestID uint; MaxQuestOrder uint;}, err error)
  FindNearbyPostInQuestWithParticularQuestOrder(lat string, lon string, radius string, questID uint, questOrder uint) (post Post, err error)
  GetNearbyCompletedPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error)
  GetNearbyUncompletedFirstPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error)
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
  databaseManager.DB.Preload("Pin").Preload("User").Where(whereClause, lat, lon, radius, questID, questOrder).First(&post)

  return
}

func (databaseManager databaseManager) GetNearbyCompletedPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error) {
  var results []struct{PostID uint}

  query := `SELECT p.id AS post_id
            FROM posts p
            INNER JOIN discovered_posts dp ON dp.user_id = ? AND dp.post_id = p.id
            WHERE ((p.latitude - ?)^2.0 + ((p.longitude - ?)* cos(p.latitude / 57.3))^2.0)  < (?^2)*0.000000000080815075
            AND dp.completed = true`

  databaseManager.DB.Raw(query, userID, lat, lon, radius).Scan(&results)

  for _, result := range results {
    var post Post
    databaseManager.DB.Preload("Pin").Preload("User").Where("id = ?", result.PostID).First(&post)
    if post.ID != 0 {
      posts = append(posts, post)
    }
  }

  return
}

func (databaseManager databaseManager) GetNearbyUncompletedFirstPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error) {
  var results []struct{PostID uint}

  query := `SELECT p.id AS post_id
            FROM posts p
            LEFT JOIN discovered_posts dp ON dp.user_id = ? AND dp.post_id = p.id
            WHERE ((p.latitude - ?)^2.0 + ((p.longitude - ?)* cos(p.latitude / 57.3))^2.0)  < (?^2)*0.000000000080815075
            AND (dp.id IS NULL OR dp.completed = false)
            AND p.quest_order = 1`

  databaseManager.DB.Raw(query, userID, lat, lon, radius).Scan(&results)

  for _, result := range results {
    var post Post
    databaseManager.DB.Preload("Pin").Preload("User").Where("id = ?", result.PostID).First(&post)
    if post.ID != 0 {
      posts = append(posts, post)
    }
  }

  return
}


package models

import (
	      "github.com/jinzhu/gorm"
        "fmt"
)

type Post struct {
	gorm.Model
  Caption string
  Latitude float64
  Longitude float64
  Location Location // This seems redundant with Latitude and Longitude above. I'm in favor of getting rid of this and keep lat/lon if it's possible.
  Pin Pin
  ImageUrl string
  ShareOnFacebook bool
  ShareOnTwitter bool
  User User
  UserID uint
  QuestID uint
  QuestOrder uint // This is the order in which the Post appears in its parent quest
  Name string
  PercentDiscovered float64 `sql:"-"`
  Completed bool `sql:"-"`
  FormattedAddress string
  Locality string
  Neighborhood string
  PostalCode string
  Country string
  AdministrativeLevelTwo string
  AdministrativeLevelOne string
  StreetNumber string
  GooglePlaceID string
  Route string
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
  if p.ShareOnFacebook || p.ShareOnTwitter {
    err = p.Share()
  }

  if err == nil {
  	err = p.SetPin()
  	err = tx.Save(p).Error
  }

  return
}

func (p *Post) SetPin() (err error) {
  _, err = NewPin(p.ImageUrl, p.ID)
  return
}

func GetPostByID(postID string) (post Post, err error){
  err = DB.Preload("Pin").Preload("User").Where("id = ?", postID).First(&post).Error
  return
}

func GetPostByUintID(postID uint) (post Post, err error){
  err = DB.Preload("Pin").Preload("User").Where("id = ?", postID).First(&post).Error
  return
}

func GetUsersPosts(userID string) (posts []Post, err error) {
  err = DB.Preload("Pin").Preload("User").Where("user_id = ?", userID).Find(&posts).Error
  return
}

func (p *Post) Share() (err error) {
  u := p.User

  if p.ShareOnFacebook {
    u.PostToFacebook(p)
  }

  if p.ShareOnTwitter {
    u.PostToTwitter(p)
  }
  
  return
}

func GetPostsNearLocationWithUserDiscoveries(lat string, lon string, userID uint, radius string, numResults int) (posts []Post, err error) {
  posts, err = GetPostsNearLocation(lat, lon, radius, numResults)

  for i, _ := range posts {
    posts[i].SetPercentDiscovered(userID)
  }

  return posts, err
}

// TODO: Should refactor to use postGIS types
func GetPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error) {
  distanceString := "((posts.latitude - " + lat + ")^2.0 + ((posts.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  whereClause := distanceString + " < (" + radius + "^2)*0.000000000080815075"
  orderClause := distanceString + " asc"
  DB.Preload("Pin").Preload("User").Where(whereClause).Order(orderClause).Limit(numResults).Find(&posts)
  return
}

func (p *Post) SetPercentDiscovered(userID uint) (err error) {
  discoveredPost := DiscoveredPost{UserID : userID, PostID : p.ID}
  err = DB.First(&discoveredPost, discoveredPost).Error; if err == nil {
    p.PercentDiscovered = discoveredPost.PercentDiscovered
    p.Completed = discoveredPost.Completed
  } else {
    p.PercentDiscovered = 0.0
    p.Completed = false
  }
  return
}

func GetNearbyPostsAndTryToDiscoverThem(user User, lat string, lon string, radius string, numPosts int) (posts []Post, err error) {
  nearbyUncompletedNonFirstPosts, _ := getNearbyUncompletedNonFirstPosts(user.ID, lat, lon, radius)
  posts = append(posts, nearbyUncompletedNonFirstPosts...)

  nearbyCompletedPosts, _ := GetNearbyCompletedPosts(user.ID, lat, lon, radius)
  posts = append(posts, nearbyCompletedPosts...)

  nearbyUncompletedFirstPosts, _ := GetNearbyUncompletedFirstPosts(user.ID, lat, lon, radius)
  posts = append(posts, nearbyUncompletedFirstPosts...)

  err = tryToDiscoverPosts(posts, &user, lat, lon); if err != nil {
    fmt.Println("ERROR: Unable to discover posts.")
  }

  for i, _ := range posts {
    posts[i].SetPercentDiscovered(user.ID)
  }

  return
}

func getNearbyUncompletedNonFirstPosts(userID uint, lat string, lon string, radius string) (tipPosts []Post, err error) {
  // This will get the quest_order and quest_id for the maximum completed post in each quest for the user.
  results, _ := GetQuestOrderForLastCompletedPostInEachQuest(userID)

  var post Post

  for _, result := range results {
    // Let's try to get the very next post in the quest (only if it's nearby). If we can't find one, then the user has completed the entire quest or the tip post is not nearby.
    post, _ = FindNearbyPostInQuestWithParticularQuestOrder(lat, lon, radius, result.QuestID, result.MaxQuestOrder + 1)
    if post.ID != 0 {
      tipPosts = append(tipPosts, post)
    }
  }

  return
}

// This is temporary until we take on the larger task of creating quests from the client app.
func (post *Post) CreateNewQuestAndSetFieldsOnPost(user *User) {
  var quest Quest
  quest.Description = post.Caption
  quest.UserID = user.ID
  DB.Create(&quest)
  fmt.Println("Created a new quest for post: ", post.ID, ".   New questID:", quest.ID)

  post.QuestID = quest.ID
  post.QuestOrder = 1
}

func CompletePostForUser(postID string, user User) (err error) {
  var discoveredPost DiscoveredPost
  err = DB.Model(&discoveredPost).Where("post_id = ? AND user_id = ?", postID, user.ID).Update("completed", true).Error
  return
}

// This will update fields that need to be updated in order for things to work. Should be removed after it's been run on all machines (dev and prod).
// Here we'll look for posts without QuestIDs. For each one we find, we'll make a quest for it.
func DatabaseUpdateTemporary() {
  var posts []Post
  DB.Find(&posts)

  for _, post := range posts {
    if post.QuestID == 0 {
      fmt.Println("Found a post with a questID of 0. post.ID: ", post.ID, " post.QuestID: ", post.QuestID)
      
      var quest Quest
      quest.Description = post.Caption
      DB.Create(&quest)
      quest.Description = "(this quest was automatically generated so things don't break)"
      fmt.Println("Created a new quest for it: ", quest.ID)
      
      post.QuestID = quest.ID
      post.QuestOrder = 1
      DB.Save(&post)
    }
  }
}

func GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error) {
  distanceString := "((posts.latitude - " + lat + ")^2.0 + ((posts.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  whereClause := distanceString + " < (" + radius + "^2)*0.000000000080815075"
  whereClause += " AND QuestOrder = 1"
  orderClause := distanceString + " asc"
  DB.Preload("Pin").Preload("User").Where(whereClause).Order(orderClause).Limit(numResults).Find(&posts)
  return
}

func GetQuestOrderForLastCompletedPostInEachQuest(userID uint) (results []struct{QuestID uint; MaxQuestOrder uint;}, err error) {
  query := `SELECT p.quest_id, MAX(quest_order) AS max_quest_order
            FROM discovered_posts dp 
            JOIN posts p ON 
              dp.user_id = ? AND 
              dp.post_id = p.id AND
              dp.Completed = true
            GROUP BY p.quest_id
            AND p.deleted_at IS NULL`

  DB.Raw(query, userID).Scan(&results)

  return
}

func FindNearbyPostInQuestWithParticularQuestOrder(lat string, lon string, radius string, questID uint, questOrder uint) (post Post, err error) {
  whereClause := `((posts.latitude - ?)^2.0 + ((posts.longitude - ?)* cos(latitude / 57.3))^2.0) < (?^2)*0.000000000080815075
                      AND quest_id = ?
                      AND quest_order = ?`
  DB.Preload("Pin").Preload("User").Where(whereClause, lat, lon, radius, questID, questOrder).First(&post)

  return
}

func GetNearbyCompletedPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error) {
  var results []struct{PostID uint}

  query := `SELECT p.id AS post_id
            FROM posts p
            INNER JOIN discovered_posts dp ON dp.user_id = ? AND dp.post_id = p.id
            WHERE ((p.latitude - ?)^2.0 + ((p.longitude - ?)* cos(p.latitude / 57.3))^2.0)  < (?^2)*0.000000000080815075
            AND dp.completed = true
            AND p.deleted_at IS NULL`

  DB.Raw(query, userID, lat, lon, radius).Scan(&results)

  for _, result := range results {
    var post Post
    DB.Preload("Pin").Preload("User").Where("id = ?", result.PostID).First(&post)
    if post.ID != 0 {
      posts = append(posts, post)
    }
  }

  return
}

func GetNearbyUncompletedFirstPosts(userID uint, lat string, lon string, radius string) (posts []Post, err error) {
  var results []struct{PostID uint}

  query := `SELECT p.id AS post_id
            FROM posts p
            LEFT JOIN discovered_posts dp ON dp.user_id = ? AND dp.post_id = p.id
            WHERE ((p.latitude - ?)^2.0 + ((p.longitude - ?)* cos(p.latitude / 57.3))^2.0)  < (?^2)*0.000000000080815075
            AND (dp.id IS NULL OR dp.completed = false)
            AND p.quest_order = 1
            AND p.deleted_at IS NULL`

  DB.Raw(query, userID, lat, lon, radius).Scan(&results)

  for _, result := range results {
    var post Post
    DB.Preload("Pin").Preload("User").Where("id = ?", result.PostID).First(&post)
    if post.ID != 0 {
      posts = append(posts, post)
    }
  }

  return
}
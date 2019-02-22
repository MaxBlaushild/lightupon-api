package models

import (
	      "github.com/jinzhu/gorm"
        "fmt"
)

type Post struct {
	gorm.Model
  Caption string
  Location Location
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
  Latitude float64
  Longitude float64
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

func GetNearbyPostsAndTryToDiscoverThem(user User, lat string, lon string, radius string, numPosts int, databaseAccessor DatabaseAccessor) (posts []Post, err error) {
  nearbyUncompletedNonFirstPosts, _ := getNearbyUncompletedNonFirstPosts(user.ID, lat, lon, radius, databaseAccessor)
  posts = append(posts, nearbyUncompletedNonFirstPosts...)

  nearbyCompletedPosts, _ := databaseAccessor.GetNearbyCompletedPosts(user.ID, lat, lon, radius)
  posts = append(posts, nearbyCompletedPosts...)

  nearbyUncompletedFirstPosts, _ := databaseAccessor.GetNearbyUncompletedFirstPosts(user.ID, lat, lon, radius)
  posts = append(posts, nearbyUncompletedFirstPosts...)

  // TODO: pass the database accessor here
  err = user.TryToDiscoverPosts(posts); if err != nil {
    fmt.Println("ERROR: Unable to discover posts.")
  }

  for i, _ := range posts {
    posts[i].SetPercentDiscovered(user.ID)
  }

  return
}

func getNearbyUncompletedNonFirstPosts(userID uint, lat string, lon string, radius string, databaseAccessor DatabaseAccessor) (tipPosts []Post, err error) {
  // This will get the quest_order and quest_id for the maximum completed post in each quest for the user.
  results, _ := databaseAccessor.GetQuestOrderForLastCompletedPostInEachQuest(userID)

  var post Post

  for _, result := range results {
    // Let's try to get the very next post in the quest (only if it's nearby). If we can't find one, then the user has completed the entire quest or the tip post is not nearby.
    post, _ = databaseAccessor.FindNearbyPostInQuestWithParticularQuestOrder(lat, lon, radius, result.QuestID, result.MaxQuestOrder + 1)
    if post.ID != 0 {
      tipPosts = append(tipPosts, post)
    }
  }

  return
}

func CompletePostForUser(postID string, user User) (err error) {
  var discoveredPost DiscoveredPost
  DB.Model(&discoveredPost).Where("post_id = ? AND user_id = ?", postID, user.ID).Update("completed", true)
  return
}

func DatabaseUpdateTemporary() {
  var post Post
  DB.Model(&post).Where("quest_id IS NULL").Update("quest_id", 1)
  DB.Model(&post).Where("quest_order IS NULL").Update("quest_order", 1)
}
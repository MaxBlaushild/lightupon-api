package models

import (
	      "github.com/jinzhu/gorm"
        "math"
        "time"
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
  Comments []Comment
  User User
  UserID uint
  Name string
  PercentDiscovered float64 `sql:"-"`
  RawScore int `sql:"-"`
  TemporalScore float64 `sql:"-"`
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
  Cost int
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

func GetPostsNearLocation(lat string, lon string, userID uint, radius string, numResults int) (posts []Post, err error) {
  distanceString := "((posts.latitude - " + lat + ")^2.0 + ((posts.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  query := distanceString + " < (" + radius + "^2)*0.000000000080815075"
  order := distanceString + " asc"
  limit := 5 * numResults
  DB.Preload("Pin").Preload("User").Where(query).Order(order).Limit(limit).Find(&posts)

  for i, _ := range posts {
    // posts[i].SetPercentDiscovered(userID)
    posts[i].PercentDiscovered = 1.0
    posts[i].SetScores()
  }

  posts = getTopNScoringPosts(posts, numResults)

  return
}

func (p *Post) SetPercentDiscovered(userID uint) (err error) {
  discoveredPost := DiscoveredPost{UserID : userID, PostID : p.ID}
  err = DB.First(&discoveredPost, discoveredPost).Error; if err == nil {
    p.PercentDiscovered = discoveredPost.PercentDiscovered
  }
  return
}

func GetRawScoreForPost(postID uint) int {
  votes := []Vote{}
  DB.Where("post_id = ?", postID).Find(&votes)
  total := 0
  for i := 0; i < len(votes); i++ {
    if votes[i].Upvote {
      total += 1
    } else {
      total += -1
    }
  }
  return total
}

// This seems like a good jumping off point for how to calculate costs for posts
func CalculateCostToPostAtLocation(lat float64, lon float64) int {
  cost := 0
  latString := fmt.Sprintf("%.6f", lat)
  lonString := fmt.Sprintf("%.6f", lon)
  posts, _ := GetPostsNearLocation(latString, lonString, 1, "100", 10)
  for i := 0; i < len(posts); i++ {
    cost = cost + GetRawScoreForPost(posts[i].ID)
  }

  if cost <= 0 {
    return 1
  }

  return cost
}

func (p *Post) SetScores() {
  p.RawScore = GetRawScoreForPost(p.ID)
  timeDiff := time.Now().Sub(p.CreatedAt).Minutes()
  p.TemporalScore = float64(p.RawScore) / math.Log(timeDiff + 1)
}

func getTopNScoringPosts(inputPosts []Post, n int) (postsToReturn []Post) {
  var topScoringIndex int
  var topScore float64
  for len(inputPosts) > 0 && len(postsToReturn) < n {
    topScore = 0; topScoringIndex = 0
    for i := 0; i < len(inputPosts); i++ {
      if inputPosts[i].TemporalScore > topScore {
        topScore = inputPosts[i].TemporalScore
        topScoringIndex = i
      }
    }
    postsToReturn = append(postsToReturn, inputPosts[topScoringIndex])
    inputPosts = removePostFromSlice(inputPosts, topScoringIndex)
  }
  return
}

func removePostFromSlice(inputPosts []Post, indexToRemove int) (postsToReturn []Post) {
  for i := 0; i < len(inputPosts); i++ {
      if i != indexToRemove {
          postsToReturn = append(postsToReturn, inputPosts[i])
      }
  }
  return
}
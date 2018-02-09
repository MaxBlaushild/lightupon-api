package models

import(
      "github.com/jinzhu/gorm"
      "lightupon-api/live"
      )

type DiscoveredPost struct {
  gorm.Model
  UserID uint
  PostID uint
  PercentDiscovered float64
}

const unlockThresholdSmall float64 = 20
const unlockThresholdLarge float64 = 200
const fadePeriod float64 = 4

func (dS *DiscoveredPost) NotFullyDiscovered() bool {
  return dS.PercentDiscovered < 1.0
}

func UpsertDiscoveredPost(discoveredPost *DiscoveredPost) {
  if DB.NewRecord(discoveredPost) {
    DB.Create(&discoveredPost)
  } else {
    DB.Model(&discoveredPost).Update("PercentDiscovered", discoveredPost.PercentDiscovered)
  }
  postUpdate := live.PostUpdate{UpdatedPostID: discoveredPost.PostID, UserID: discoveredPost.UserID}
  live.Hub.UpdateClient <- postUpdate
}

func (dS *DiscoveredPost) UpdatePercentDiscovered(user *User, post *Post) {
  newPercentDiscovered := calculatePercentDiscovered(user, post)
  if (newPercentDiscovered > dS.PercentDiscovered) {
    dS.PercentDiscovered = newPercentDiscovered
    UpsertDiscoveredPost(dS)
  }
}

func calculatePercentDiscovered(user *User, post *Post) (percentDiscovered float64) {
  distance := CalculateDistance(user.Location, post.Location)
  if (distance < unlockThresholdSmall) {
    percentDiscovered = 1.0
  } else if (distance > unlockThresholdLarge) {
    percentDiscovered = 0.0
  } else {
    percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall)) // TODO: Update this to be a nice smoove cosine function
  }
  return
}

func GetCurrentDiscoveredPost(userID uint, postID uint) DiscoveredPost {
  discoveredPost := DiscoveredPost{UserID: userID, PostID: postID}
  DB.First(&discoveredPost, discoveredPost)
  return discoveredPost
}
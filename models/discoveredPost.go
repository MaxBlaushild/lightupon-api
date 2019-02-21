package models

import(
      "github.com/jinzhu/gorm"
      )

type DiscoveredPost struct {
  gorm.Model
  UserID uint
  PostID uint
  PercentDiscovered float64
  Completed bool
}

const unlockThresholdSmall float64 = 10
const unlockThresholdLarge float64 = 40

func saveNewPercentDiscoveredToDB(user *User, post *Post, newPercentDiscovered float64) {
  discoveredPost := getDiscoveredPostOrCreateNew(user.ID, post.ID)
  DB.Model(&discoveredPost).Update("PercentDiscovered", newPercentDiscovered)
}

func tryToDiscoverPost(post *Post, user *User) {
  if post.PercentDiscovered == 1.0 {
    return
  }

  newPercentDiscovered := calculatePercentDiscovered(user, post)

  if (newPercentDiscovered > post.PercentDiscovered) {
    saveNewPercentDiscoveredToDB(user, post, newPercentDiscovered)
  }

  return
}

func calculatePercentDiscovered(user *User, post *Post) (percentDiscovered float64) {
  distance := CalculateDistance(user.Location, Location{Latitude: post.Latitude, Longitude: post.Longitude})
  if (distance < unlockThresholdSmall) {
    percentDiscovered = 1.0
  } else if (distance > unlockThresholdLarge) {
    percentDiscovered = 0.0
  } else {
    percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall))
  }
  return
}

func getDiscoveredPostOrCreateNew(userID uint, postID uint) DiscoveredPost {
  discoveredPost := DiscoveredPost{UserID: userID, PostID: postID}
  DB.First(&discoveredPost, discoveredPost)
  if discoveredPost.ID == 0 {
    DB.Create(&discoveredPost)
  }
  return discoveredPost
}

func getNearbyDiscoveredPosts(userID uint, postID uint) DiscoveredPost {
  discoveredPost := DiscoveredPost{UserID: userID, PostID: postID}
  DB.First(&discoveredPost, discoveredPost)
  if discoveredPost.ID == 0 {
    DB.Create(&discoveredPost)
  }
  return discoveredPost
}
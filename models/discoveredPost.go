package models

import(
      "github.com/jinzhu/gorm"
      )

type DiscoveredPost struct {
  gorm.Model
  UserID uint
  PostID uint
  PercentDiscovered float64
}

const unlockThresholdSmall float64 = 10
const unlockThresholdLarge float64 = 40

func saveNewPercentDiscoveredToDB(user *User, post *Post, newPercentDiscovered float64) {
  discoveredPost := GetDiscoveredPostOrCreateNew(user.ID, post.ID)
  DB.Model(&discoveredPost).Update("PercentDiscovered", newPercentDiscovered)
}

func tryToDiscover(post *Post, user *User) {
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

func GetDiscoveredPostOrCreateNew(userID uint, postID uint) DiscoveredPost {
  discoveredPost := DiscoveredPost{UserID: userID, PostID: postID}
  DB.First(&discoveredPost, discoveredPost)
  if discoveredPost.ID == 0 {
    DB.Create(&discoveredPost)
  }
  return discoveredPost
}
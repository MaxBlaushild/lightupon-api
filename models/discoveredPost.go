package models

import(
      "github.com/jinzhu/gorm"
      "strconv"
      )

type DiscoveredPost struct {
  gorm.Model
  UserID uint
  PostID uint
  Post Post
  PercentDiscovered float64
  Completed bool
}

const unlockThresholdSmall float64 = 20
const unlockThresholdLarge float64 = 80

func tryToDiscoverPosts(posts []Post, user *User, lat string, lon string) (err error)  {
  for _, post := range posts {
    if !post.Completed {
      newPercentDiscovered, completed := calculatePercentDiscovered(&post, lat, lon)

      if (newPercentDiscovered > post.PercentDiscovered) || completed {
        discoveredPost := getDiscoveredPostOrCreateNew(user.ID, post.ID)
        DB.Model(&discoveredPost).Update("PercentDiscovered", newPercentDiscovered)
      }
    }
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

func calculatePercentDiscovered(post *Post, lat string, lon string) (percentDiscovered float64, completed bool) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)

  distance := CalculateDistance(Location{Latitude: latFloat, Longitude: lonFloat}, Location{Latitude: post.Latitude, Longitude: post.Longitude})
  if (distance < unlockThresholdSmall) {
    percentDiscovered = 1.0
    completed = true
  } else if (distance > unlockThresholdLarge) {
    percentDiscovered = 0.0
    completed = false
  } else {
    percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall))
    completed = false
  }
  return
}
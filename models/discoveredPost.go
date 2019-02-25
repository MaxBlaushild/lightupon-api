package models

import(
      "github.com/jinzhu/gorm"
      "fmt"
      "strconv"
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

func tryToDiscoverPosts(posts []Post, user *User, lat string, lon string) (err error)  {
  for i, _ := range posts {
    tryToDiscoverPost(&posts[i], user, lat, lon)
  }

  return
}

func tryToDiscoverPost(post *Post, user *User, lat string, lon string) {
  fmt.Println("____________trying to discover post: ", post.ID, " for user: ", user.ID)
  if post.PercentDiscovered == 1.0 {
    fmt.Println("    fully discovered already.")
    return
  }

  newPercentDiscovered := calculatePercentDiscovered(post, lat, lon)
  fmt.Println("    newPercentDiscovered: ", newPercentDiscovered)

  if (newPercentDiscovered > post.PercentDiscovered) {
    saveNewPercentDiscoveredToDB(user, post, newPercentDiscovered)
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

func calculatePercentDiscovered(post *Post, lat string, lon string) (percentDiscovered float64) {
  latFloat, _ := strconv.ParseFloat(lat, 64)
  lonFloat, _ := strconv.ParseFloat(lon, 64)

  distance := CalculateDistance(Location{Latitude: latFloat, Longitude: lonFloat}, Location{Latitude: post.Latitude, Longitude: post.Longitude})
  if (distance < unlockThresholdSmall) {
    percentDiscovered = 1.0
  } else if (distance > unlockThresholdLarge) {
    percentDiscovered = 0.0
  } else {
    percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall))
  }
  return
}

func saveNewPercentDiscoveredToDB(user *User, post *Post, newPercentDiscovered float64) {
  discoveredPost := getDiscoveredPostOrCreateNew(user.ID, post.ID)
  DB.Model(&discoveredPost).Update("PercentDiscovered", newPercentDiscovered)
}
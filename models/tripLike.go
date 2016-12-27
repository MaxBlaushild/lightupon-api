package models

import(
      "github.com/jinzhu/gorm"
      // "github.com/davecgh/go-spew/spew"
      // "fmt"
      )

type TripLike struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  TripID uint `gorm:"not null"`
}

func GetTripLikesForUser(userID uint) (likes []TripLike){
  DB.Where("user_id = $1", userID).Find(&likes)
  return
}

func GetLikesForTrip(tripID uint) (likes []TripLike){
  DB.Where("trip_id = $1", tripID).Find(&likes)
  return
}

func HasUserLikedTrip(userID uint, tripID uint) bool {
  like := TripLike{}
  DB.Where("trip_id = $1 AND user_id = $2", tripID, userID).First(&like)
  
  // If the like's ID is zero, then nothing was found in the DB for that trip/user combo
  if like.ID == 0 {
    return false
  } else {
    return true
  }
}

func GetTotalLikesForTrip(tripID uint) int {
  likes := []TripLike{}
  DB.Where("trip_id = $1", tripID).Find(&likes)

  return len(likes)
}
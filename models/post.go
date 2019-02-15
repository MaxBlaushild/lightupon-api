package models

import (
	      "github.com/jinzhu/gorm"
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
  Name string
  PercentDiscovered float64 `sql:"-"`
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

// TODO: Should refactor to use postGIS types but GORM doesn't support them, so that's a larger discussion
func GetPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []Post, err error) {
  distanceString := "((posts.latitude - " + lat + ")^2.0 + ((posts.longitude - " + lon + ")* cos(latitude / 57.3))^2.0)"
  query := distanceString + " < (" + radius + "^2)*0.000000000080815075"
  order := distanceString + " asc"
  limit := 5 * numResults
  DB.Preload("Pin").Preload("User").Where(query).Order(order).Limit(limit).Find(&posts) // Why are we preloading the user here? Does it matter who created the post?

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
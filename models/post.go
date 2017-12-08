package models

import (
	      "github.com/jinzhu/gorm"
        "fmt"
)

type Post struct {
	gorm.Model
  Caption string
  Location Location
  Pin Pin `gorm:"polymorphic:Owner;"`
  ImageUrl string
  ShareOnFacebook bool
  ShareOnTwitter bool
  Comments []Comment
  User User
  UserID uint
  Name string
  PercentDiscovered float64 `sql:"-"`
  RawScore float64 `sql:"-"`
  TimeVoteScore float64 `sql:"-"`
  SpatialScore float64 `sql:"-"`
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

func (p *Post) SetPin() error {
  fmt.Println("trying to create pin")
  _, err := NewPin(p.ImageUrl, p.ID, "Post")
  return err
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
  limit := 3 * numResults
  DB.Preload("Pin").Where(query).Order(order).Limit(limit).Find(&posts)

  for i, _ := range posts {
    posts[i].SetPercentDiscovered(userID)
  }

  return
}

func (p *Post) SetPercentDiscovered(userID uint) (err error) {
  discoveredPost := DiscoveredPost{UserID : userID, PostID : p.ID}
  err = DB.First(&discoveredPost, discoveredPost).Error; if err == nil {
    p.PercentDiscovered = discoveredPost.PercentDiscovered
  }
  return
}
package models

import(
      "github.com/jinzhu/gorm"
      )

type Card struct {
  gorm.Model
  Caption string
  Latitude float64
  Longitude float64
  ImageURL string
  Pin string
  SelectedPin string
  SceneID uint
  ShareOnFacebook bool
  Comments []Comment
  CardOrder uint
  Universal bool
  NibID string `gorm:"not null"`
}

func (c *Card) AfterCreate(tx *gorm.DB) (err error) {
  if c.ShareOnFacebook {
    err = c.Share()
  }
  return
}

func (c *Card) Share() (err error) {
  u, err := c.User()

  if c.ShareOnFacebook {
    u.PostToFacebook(c)
  }
  
  return
}

func (c *Card) User() (user User, err error) {
  scene := Scene{}
  trip := Trip{}
  err = DB.First(&scene, c.SceneID).Error
  err = DB.First(&trip, scene.TripID).Error
  err = DB.First(&user, trip.UserID).Error
  return
}

func ShiftCardsUp(cardOrder int, sceneID int) bool {
  card := Card{}
  DB.Where("scene_id = $1 AND card_order = $2", sceneID, cardOrder).First(&card)
  if card.ID == 0 {
    return true
  } else {
    ShiftCardsUp(cardOrder + 1, 1)
    DB.Model(&card).Update("card_order", cardOrder + 1)
    return true
  }
}

func ShiftCardsDown(cardOrder int, sceneID int) bool {
  card := Card{}
  DB.Where("scene_id = $1 AND card_order = $2", sceneID, cardOrder + 1).First(&card)
  if card.ID == 0 {
    return true
  } else {
    ShiftCardsUp(cardOrder + 1, 1)
    DB.Model(&card).Update("card_order", cardOrder)
    return true
  }
}

package models

import(
      "github.com/jinzhu/gorm"
      "fmt"
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
  UserID uint
  ShareOnFacebook bool
  Comments []Comment
  CardOrder uint
  Universal bool
  NibID string `gorm:"not null"`
}

func (c *Card) AfterCreate(tx *gorm.DB) (err error) {
  fmt.Println("in after create hook:")
  fmt.Println(c.ShareOnFacebook)
  if c.ShareOnFacebook {
    fmt.Println("sharing")
    err = c.Share()
  }
  return
}

func (c *Card) Share() (err error) {
  fmt.Println("in share")
  u := User{}
  DB.First(&u, c.UserID)
  fmt.Println("found user")
  fmt.Println(u)
  if c.ShareOnFacebook {
    fmt.Println("about to share to facebook")
    u.PostToFacebook(c)
  }
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

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
  SceneID uint
  Comments []Comment
  CardOrder uint
  Universal bool
  NibID string `gorm:"not null"`
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


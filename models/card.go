package models

import(
      "github.com/jinzhu/gorm"
      )

type Card struct {
  gorm.Model
  CardType string
  Text string
  ImageURL string
  SceneID uint
  CardOrder uint
  Universal bool
}
package models

import(
      "github.com/jinzhu/gorm"
      )

type Card struct {
  gorm.Model
  Dialogue string
  SceneID uint
  CardOrder uint
  Universal bool
}
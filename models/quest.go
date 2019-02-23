package models

import (
	      "github.com/jinzhu/gorm"
)

type Quest struct {
	gorm.Model
	Description string
}
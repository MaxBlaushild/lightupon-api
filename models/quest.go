package models

import (
	      "github.com/jinzhu/gorm"
)

type Quest struct {
	gorm.Model
	Description string
	TimeToComplete int // Measured in minutes
	UserID uint
	QuestProgress QuestProgress `sql:"-"`
	Posts []Post
}
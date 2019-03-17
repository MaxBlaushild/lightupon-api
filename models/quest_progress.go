package models

import "github.com/jinzhu/gorm"

type QuestProgress struct {
	gorm.Model
	CompletedPosts int `gorm:"default:0"`
	QuestID uint 
	UserID uint
	User User
	IsCompleted bool `sql:"not null"`
}
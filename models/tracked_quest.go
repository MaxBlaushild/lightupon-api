package models

import "github.com/jinzhu/gorm"

type TrackedQuest struct {
	gorm.Model
	QuestID uint
	Quest Quest
	User User
	UserID uint
}
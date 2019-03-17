package models

import (
	      "github.com/jinzhu/gorm"
)

type Quest struct {
	gorm.Model
	Description string
	TimeToComplete int // Measured in minutes
	UserID uint
	QuestProgress QuestProgress
	QuestProgressID uint
	Posts []Post
}

func (q *Quest) IsFinished () bool {
	return q.QuestProgress.CompletedPosts >= len(q.Posts)
}

func GetQuestWithUserContext(questID uint, userID uint) (quest Quest, err error) {
	err = DB.Preload("Posts").First(&quest, questID).Error

	if err != nil {
		return
	}

	err = DB.Where("user_id = ? and quest_id =?", userID, questID).FirstOrCreate(&quest.QuestProgress).Error
	return
}
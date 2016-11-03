package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	FollowingUserID uint `gorm:"unique_index:idx_followers"`
	FollowedUserID uint `gorm:"unique_index:idx_followers"`
}
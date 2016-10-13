package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	FollowingUser uint
	FollowedUser uint
}
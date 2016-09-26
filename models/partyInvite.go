package models

import(
      "github.com/jinzhu/gorm"
      )

type PartyInvite struct {
  gorm.Model
  PartyID uint
  UserID uint
  New bool `gorm:"default:true"`
}

// I guess there are all kinds of properties for invites that could be useful like Accepted, Seen or Declined,
// but for now I think the most important one is New because we only need to decide whether or not to alert the user of the invite

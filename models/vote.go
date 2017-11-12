package models

import(
      "github.com/jinzhu/gorm"
      // "github.com/davecgh/go-spew/spew"
      "fmt"
      "errors"
      )

type Vote struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  SceneID uint
  Upvote bool // true for upvote, false for downvote
}

func SaveVote(userID uint, sceneID uint, upvote bool) error {
	if voteHasAlreadyBeenCast(userID, sceneID) {
		return errors.New("user has already voted for this scene")
	}
	vote := Vote{UserID: userID, SceneID: sceneID, Upvote: upvote}
	DB.Create(&vote)
	return nil
}

func DeleteVote(userID uint, sceneID uint) error {
	if !voteHasAlreadyBeenCast(userID, sceneID) {
		return errors.New("user has not voted for this scene")
	}
	vote := Vote{UserID: userID, SceneID: sceneID}
	DB.Delete(&vote)
	return nil
}

func GetVoteTotalForScene(sceneID uint) int {
	votes := []Vote{}
	DB.Where("scene_id = ?", sceneID).Find(&votes)
	fmt.Println("votes", votes)
	total := 0
	for i := 0; i < len(votes); i++ {
		if votes[i].Upvote {
			total += 1
		} else {
			total += -1
		}
	}
	return total
}

func voteHasAlreadyBeenCast(userID uint, sceneID uint) bool {
	vote := Vote{UserID : userID, SceneID : sceneID}
	testVote := Vote{}
	DB.Where(&vote).First(&testVote)
	if testVote.ID != 0 {
		return true
	}
	return false
}

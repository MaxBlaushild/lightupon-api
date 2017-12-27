package models

import(
      "github.com/jinzhu/gorm"
      "errors"
      "math"
      )

type Vote struct {
  gorm.Model
  UserID uint `gorm:"not null"`
  PostID uint
  Upvote bool // true for upvote, false for downvote
}

func SaveVote(userID uint, postID uint, upvote bool) error {
	if voteHasAlreadyBeenCast(userID, postID) {
		return errors.New("user has already voted for this post")
	}
	vote := Vote{UserID: userID, PostID: postID, Upvote: upvote}
	DB.Create(&vote)
	return nil
}

func DeleteVote(userID uint, postID uint) error {
	if !voteHasAlreadyBeenCast(userID, postID) {
		return errors.New("user has not voted for this scene")
	}
	vote := Vote{UserID: userID, PostID: postID}
	DB.Delete(&vote)
	return nil
}

func voteHasAlreadyBeenCast(userID uint, postID uint) bool {
	vote := Vote{UserID : userID, PostID : postID}
	testVote := Vote{}
	DB.Where(&vote).First(&testVote)
	if testVote.ID != 0 {
		return true
	}
	return false
}

func GetManaTotalForUser(userID uint) int {
	manaTotal := 100 // i guess let's start everybody at 100 mana for now
	manaTotal = manaTotal + getManaTotalForSubmittedPosts(userID)
	manaTotal = manaTotal + getManaTotalForSubmittedVotes(userID)
	return manaTotal
}

func getManaTotalForSubmittedPosts(userID uint) int {
	manaTotal := 0
	var posts []Post
	DB.Where("user_id = ?", userID).Find(&posts)
	for i := 0; i < len(posts); i++ {
		manaTotal = manaTotal - posts[i].Cost + GetRawScoreForPost(posts[i].ID)
	}
	return manaTotal
}

func getManaTotalForSubmittedVotes(userID uint) int {
	votes := []Vote{}
	DB.Where("user_id = ?", userID).Find(&votes)
	return len(votes) //* calculateVoteEffectiveness(votes)
}

// don't know about this yet... maybe we'll need this to deal with weird incentives created by rewarding users for giving out votes
func calculateVoteEffectiveness(votes []Vote) float64 {
	upvotes := 0
	downvotes := 0
	for i := 0; i < len(votes); i++ {
		if votes[i].Upvote {
			upvotes = upvotes + 1
		} else {
			downvotes = downvotes + 1
		}
	}
	return 1.0 - math.Abs((float64(upvotes) / float64(upvotes + downvotes)) - 0.5) // should equal 1 if upvotes = downvotes and will go down to 0.5 if there is an imbalance
}
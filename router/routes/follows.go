package routes

import(
       "net/http"
       "lightupon-api/models"
)

func FollowHandler(w http.ResponseWriter, r *http.Request) {
  followingUser := GetUserFromRequest(r)
  userToFollow := GetUIntFromVars(r, "userID")

  follow := models.Follow{FollowingUserID:followingUser.ID, FollowedUserID:userToFollow}

  models.DB.FirstOrCreate(&models.Follow{}, &follow) // using FirstOrCreate just to avoid adding dupes to the DB

  respondWithCreated(w, "You just followed the shit out of that user!")
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
  followingUser := GetUserFromRequest(r)
  userToUnfollow := GetUIntFromVars(r, "userID")

  follow := models.Follow{FollowingUserID:followingUser.ID, FollowedUserID:userToUnfollow}
  models.DB.Where(follow).Delete(&models.Follow{})

  respondWithNoContent(w, "You just unfollowed the shit out of that user!")
}
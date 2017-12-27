package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/mux"
       "strconv"
       "encoding/json"
       )

func PostUpvoteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  postIDint, _ := strconv.Atoi(vars["postID"])
  postID := uint(postIDint)
  if models.SaveVote(user.ID, postID, true) == nil {
  	respondWithCreated(w, "vote was saved")
  } else {
  	respondeWithRecordExists(w, "vote has already been cast")
  }
}

func PostDownvoteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  postIDint, _ := strconv.Atoi(vars["postID"])
  postID := uint(postIDint)
  if models.SaveVote(user.ID, postID, false) == nil {
  	respondWithCreated(w, "vote was saved")
  } else {
  	respondeWithRecordExists(w, "vote has already been cast")
  }
}

func DeleteVoteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  postIDint, _ := strconv.Atoi(vars["postID"])
  postID := uint(postIDint)
  if models.DeleteVote(user.ID, postID) == nil {
  	respondWithCreated(w, "vote was deleted")
  } else {
  	respondWithNotFound(w, "vote does not exist")
  }
}

func GetRawScoreHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  postIDint, _ := strconv.Atoi(vars["postID"])
  json.NewEncoder(w).Encode(struct {VoteTotal int} { models.GetRawScoreForPost(uint(postIDint)) })
}
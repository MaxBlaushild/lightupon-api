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
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  if models.SaveVote(user.ID, sceneID, true) == nil {
  	respondWithCreated(w, "vote was saved")
  } else {
  	respondeWithRecordExists(w, "vote has already been cast")
  }
}

func PostDownvoteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  if models.SaveVote(user.ID, sceneID, false) == nil {
  	respondWithCreated(w, "vote was saved")
  } else {
  	respondeWithRecordExists(w, "vote has already been cast")
  }
}

func DeleteVoteHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  if models.DeleteVote(user.ID, sceneID) == nil {
  	respondWithCreated(w, "vote was deleted")
  } else {
  	respondWithNotFound(w, "vote does not exist")
  }
}

func GetVoteTotalHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  json.NewEncoder(w).Encode(struct {VoteTotal int} { models.GetVoteTotalForScene(uint(sceneIDint)) })
}

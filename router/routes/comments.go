package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func TripCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  tripIDString := vars["tripID"]
  tripID, err := strconv.Atoi(tripIDString)

  if err != nil {
  	respondWithBadRequest(w, "The trip ID you sent us was bunk.")
  } else {
  	 comments := models.GetCommentsForTrip(tripID)
  	json.NewEncoder(w).Encode(comments)
  }

}

func SceneCommentsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDString := vars["sceneID"]
  sceneID, err := strconv.Atoi(sceneIDString)
  if err != nil {
    respondWithBadRequest(w, "The scene ID you sent us was bunk.")
  } else {
     comments := models.GetCommentsForScene(sceneID)
    json.NewEncoder(w).Encode(comments)
  }
}

func CardCommentsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cardIDString := vars["cardID"]
  cardID, err := strconv.Atoi(cardIDString)

  if err != nil {
    respondWithBadRequest(w, "The card ID you sent us was bunk.")
  } else {
    comments := models.GetCommentsForCard(cardID)
    json.NewEncoder(w).Encode(comments)
  }

}

func PostTripCommentHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  tripIDString := vars["tripID"]
  tripID, _ := strconv.Atoi(tripIDString)

  comment := models.Comment{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&comment); if err != nil {
    respondWithBadRequest(w, "The comment you sent us was bunk.")
  }

  user := GetUserFromRequest(r)

  comment.TripID = uint(tripID)
  comment.UserID = user.ID

  models.DB.Create(&comment)

}

func PostSceneCommentHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDString := vars["sceneID"]
  sceneID, _ := strconv.Atoi(sceneIDString)

  comment := models.Comment{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&comment); if err != nil {
    respondWithBadRequest(w, "The comment you sent us was bunk.")
  }

  user := GetUserFromRequest(r)

  comment.SceneID = uint(sceneID)
  comment.UserID = user.ID

  models.DB.Create(&comment)


}

func PostCardCommentHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cardIDString := vars["cardID"]
  cardID, _ := strconv.Atoi(cardIDString)

  comment := models.Comment{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&comment); if err != nil {
    respondWithBadRequest(w, "The comment you sent us was bunk.")
  }

  user := GetUserFromRequest(r)

  comment.CardID = uint(cardID)
  comment.UserID = user.ID

  models.DB.Create(&comment)

}
package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )


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

package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/mux"
       "strconv"
       )

func LikeSceneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := GetUserFromRequest(r)
  id, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(id)
  sceneLike := models.SceneLike{User: user, SceneID: sceneID}
  models.DB.Create(&sceneLike)
  respondWithCreated(w, "The scene was liked.")
}

func UnlikeSceneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := GetUserFromRequest(r)
  id, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(id)
  sceneLike := models.SceneLike{UserID: user.ID, SceneID: sceneID}
  models.DB.Where(&sceneLike).Delete(models.SceneLike{})
  respondWithCreated(w, "The scene was unliked.")
}
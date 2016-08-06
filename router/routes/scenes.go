package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func PopularScenesHandler(w http.ResponseWriter, r *http.Request) {
  scenes := []models.Scene{}
  models.DB.Where("Featured = true").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

func ScenesHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["tripId"])
  scenes := []models.Scene{}
  models.DB.Find(models.Trip{}, id).Association("Scenes").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

// request should look like {"SceneOrder":3, "Name":"new scene", "Latitude":76.567,"Longitude":87.345}
func CreateSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  tripID, _ := strconv.Atoi(vars["tripID"])
  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&scene)
  if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk.")
  }

  scene.TripID = uint(tripID)
  if (scene.ID > 0) {
    // If wants to use a ole ass scene for their trip, so let's clone it and reset the ID to zero
    newSceneOrder := scene.SceneOrder // Grab the requested new SceneOrder because reflecting on the next line will destroy it
    models.DB.Find(&scene)
    scene.ID = 0 // Set the sceneID to zero so it will insert properly below
    scene.SceneOrder = newSceneOrder
    scene.Featured = false
  }
  models.ShiftScenesUp(int(scene.SceneOrder), tripID)
  models.DB.Create(&scene)
  respondWithCreated(w, "The scene was created.")
}

func DeleteSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  scene := models.Scene{}
  scene.ID = sceneID
  models.DB.Find(&scene)
  models.ShiftScenesDown(int(scene.SceneOrder), int(scene.TripID))
  models.DB.Delete(&scene)
  respondWithNoContent(w, "The scene was deleted.")
}
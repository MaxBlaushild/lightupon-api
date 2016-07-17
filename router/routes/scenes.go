package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "fmt"
       )

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
  if err != nil {fmt.Println(err)}

  scene.TripID = uint(tripID)

  // models.DB.Query("UPDATE scenes SET scene_order = scene_order + 1 WHERE trip_id = $1 AND scene_order > $2", tripID, newSceneOrder)

  models.DB.Create(&scene)
}
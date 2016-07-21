package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "fmt"
       // "github.com/davecgh/go-spew/spew"
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
  models.ShiftScenesUp(int(scene.SceneOrder), tripID)
  scene.TripID = uint(tripID)
  models.DB.Create(&scene)
}

func DeleteSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  scene := models.Scene{}
  scene.ID = sceneID
  models.DB.Delete(&scene)
}

func DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cardIDint, _ := strconv.Atoi(vars["cardID"])
  cardID := uint(cardIDint)
  card := models.Card{}
  card.ID = cardID
  models.DB.Delete(&card)
}

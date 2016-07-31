package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "github.com/davecgh/go-spew/spew"
       "fmt"
       )

func PopularScenesHandler(w http.ResponseWriter, r *http.Request) {
  scenes := []models.Scene{}
  models.DB.Find(&scenes)
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
  fmt.Println("spew scene")
  spew.Dump(scene)
  if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk.")
  }
  scene.TripID = uint(tripID)
  // TODO: add logic to say "if the sceneID is not null, then the scene already exists and we just need to add it to the trip"
  if (scene.ID > 0) {
    fmt.Print("ok yeah we can't do that just yet")  
  } else {
    models.ShiftScenesUp(int(scene.SceneOrder), tripID)
    models.DB.Create(&scene)
    respondWithCreated(w, "The scene was created.")
  }
}

func DeleteSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  scene := models.Scene{}
  scene.ID = sceneID
  models.DB.Delete(&scene)
  respondWithNoContent(w, "The scene was deleted.")
}

func DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  cardIDint, _ := strconv.Atoi(vars["cardID"])
  cardID := uint(cardIDint)
  card := models.Card{}
  card.ID = cardID
  models.DB.Delete(&card)
  models.ShiftCardsDown(int(card.CardOrder), int(card.SceneID))
  respondWithNoContent(w, "The card was deleted.")
}

func DeleteTripHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  tripIDint, _ := strconv.Atoi(vars["tripID"])
  tripID := uint(tripIDint)
  trip := models.Trip{}
  trip.ID = tripID
  models.DB.Delete(&trip)
  respondWithNoContent(w, "The trip was deleted.")
}



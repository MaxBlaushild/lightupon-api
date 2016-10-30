package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "github.com/jinzhu/gorm"
       "lightupon-api/googleMaps"
       )

func TripsHandler(w http.ResponseWriter, r *http.Request) {
  lat, lon := GetUserLocationFromRequest(r)
  trips := []models.Trip{}
  models.DB.Preload("Locations").Preload("User").Preload("Scenes", func(DB *gorm.DB) *gorm.DB {
    return DB.Order("Scenes.scene_order ASC") // Preload and order scenes for the map view
  }).Order("((latitude - " + lat + ")^2.0 + ((longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc;").Find(&trips)

  for i := 0; i < len(trips); i++ {
    snappedLocations, err := googleMaps.SnapLocations(trips[i].Locations); if err == nil {
      trips[i].PutLocations(snappedLocations)
    }
  }

  json.NewEncoder(w).Encode(trips)
}

func CreateSelfieTripHandler(w http.ResponseWriter, r *http.Request) {
  selfie := models.Selfie{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&selfie)
  if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk.")
    return
  }

  trip := models.Trip{Title: "New Selfie" }
  scene := models.Scene{ 
    Latitude: selfie.Location.Latitude, 
    Longitude: selfie.Location.Longitude, 
    SceneOrder: 1, 
    BackgroundUrl: selfie.ImageUrl,
  }
  card := models.Card{ NibID: "PictureHero", ImageURL: selfie.ImageUrl }
  scene.Cards = append (scene.Cards, card)
  trip.Scenes = append (trip.Scenes, scene)
  models.DB.Create(&trip)
  respondWithCreated(w, "The trip was created.")
}

func TripHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
  trip := models.Trip{}
  models.DB.First(&trip, id)
  json.NewEncoder(w).Encode(trip)
}

func CreateTripHandler(w http.ResponseWriter, r *http.Request) {
  // request body should consist of {Title: "Balls"}
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip)
  if err != nil {
    respondWithBadRequest(w, "The trip credentials you sent us were wack!")
  }

  user := GetUserFromRequest(r)
  trip.UserID = user.ID
  models.DB.Create(&trip)
  json.NewEncoder(w).Encode(trip)
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

func GetTripsForUserHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  trips := []models.Trip{}
  models.DB.Where("owner = $1", user.ID).Find(&trips)
  json.NewEncoder(w).Encode(trips)
}

func AdminCreateTripHandler(w http.ResponseWriter, r *http.Request) {
  // request body should consist of {Title: "Balls"}
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip)
  if err != nil {
    respondWithBadRequest(w, "The trip credentials you sent us were wack!")
  }
  trip.UserID = 1
  models.DB.Create(&trip)
  json.NewEncoder(w).Encode(trip)
}
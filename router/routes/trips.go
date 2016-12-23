package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "fmt"
       )

func TripsHandler(w http.ResponseWriter, r *http.Request) {  
  lat, lon := GetUserLocationFromRequest(r)
  trips := models.GetTripsNearLocation(lat, lon)

  json.NewEncoder(w).Encode(trips)
}

func CreateDegenerateTripHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  fmt.Println("INFO: Creating selfie trip")
  selfie := models.Selfie{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&selfie)
  if err != nil {
    respondWithBadRequest(w, "The selfie you sent us was bunk.")
    return
  }

  models.CreateSelfieTrip(selfie, user.ID)

  respondWithCreated(w, "The trip was created.")
}

func TripHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r) 
  tripID, err := strconv.Atoi(vars["id"]); if err != nil {
    respondWithBadRequest(w, "The id you sent us was bunk.")
    return
  }
  
  trip := models.GetTrip(tripID)
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
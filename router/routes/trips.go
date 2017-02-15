package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "fmt"
       )

func NearbyTripsHandler(w http.ResponseWriter, r *http.Request) {  
  user := GetUserFromRequest(r)

  lat, lon := GetUserLocationFromRequest(r)
  trips := models.GetTripsNearLocation(lat, lon, user.ID)

  json.NewEncoder(w).Encode(trips)
}

func UpdateActiveTrip(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip); if err != nil {
    respondWithBadRequest(w, "The trip you sent us was bunk.")
  }
  user := GetUserFromRequest(r)
  activeTrip := user.ActiveTrip()
  models.DB.Model(&activeTrip).Update(trip)
  json.NewEncoder(w).Encode(activeTrip)
}

func ActiveTripHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeTrip := user.ActiveTrip()
  json.NewEncoder(w).Encode(activeTrip)

}

func GetUsersTripsHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  userID := vars["userID"]
  trips := models.GetTripsForUser(userID)
  json.NewEncoder(w).Encode(trips)
}

func TripHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r) 
  tripID, err := strconv.Atoi(vars["id"]); if err != nil {
    respondWithBadRequest(w, "The id you sent us was bunk.")
    return
  }

  user := GetUserFromRequest(r)

  trip := models.GetTrip(tripID, user.ID)

  json.NewEncoder(w).Encode(trip)
}

func CreateDegenerateTripHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&scene); if err != nil {
    respondWithBadRequest(w, "The selfie you sent us was bunk.")
    return
  }

  trip := models.Trip{
    Active: false,
    Scenes: []models.Scene{scene},
  }

  user.AddTrip(&trip)

  respondWithCreated(w, "The trip was created.")
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

func LikeTripHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  vars := mux.Vars(r)
  tripIDint, _ := strconv.Atoi(vars["tripID"])
  tripID := uint(tripIDint)

  like := models.TripLike{}  // why the fuck can't I get struct literal declarations to work

  // Check if user has already liked this trip
  if (models.HasUserLikedTrip(user.ID, tripID)) {
    // This shoudn't happen, so maybe log an error in the future if we feel like it
    fmt.Println("user tried to like a trip twice")
  } else {
    like.UserID = user.ID
    like.TripID = tripID

    models.DB.Create(&like)
  }

  return
}
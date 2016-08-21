package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "github.com/jinzhu/gorm"
       )

func TripsHandler(w http.ResponseWriter, r *http.Request) {
  lat, lon := GetUserLocationFromRequest(r)
  latString := strconv.FormatFloat(lat, 'f', 6, 64)
  lonString := strconv.FormatFloat(lon, 'f', 6, 64)

  trips := []models.Trip{}
  models.DB.Preload("Scenes", func(DB *gorm.DB) *gorm.DB {
    return DB.Order("Scenes.scene_order ASC") // Preload and order scenes for the map view
  }).Order("((latitude - " + latString + ")^2.0 + ((longitude - " + lonString + ")* cos(latitude / 57.3))^2.0) asc;").Find(&trips)

  json.NewEncoder(w).Encode(trips)
}

// To be used with something like http://www.localhost:5000/nearby_trips?lat=42.33865&lon=-71.079994&threshold=1
// In this example, it will return all trips within one mile of headquarters
func NearbyTripsHandler(w http.ResponseWriter, r *http.Request) {
  // Pull all necessary parameters out of the URL
  query := r.URL.Query()
  lat, _ := strconv.ParseFloat(query["lat"][0], 64)
  lon, _ := strconv.ParseFloat(query["lon"][0], 64)
  threshold, _ := strconv.ParseFloat(query["threshold"][0], 64)

  trips := []models.Trip{}
  models.DB.Where("(pow(latitude - $1, 2) + pow(longitude - $2,2)) < $3 ORDER BY title", lat, lon, threshold/10000).Find(&trips)
  json.NewEncoder(w).Encode(trips)
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
  trip.Owner = int(user.ID)
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
  trip.Owner = 1
  models.DB.Create(&trip)
  json.NewEncoder(w).Encode(trip)
}
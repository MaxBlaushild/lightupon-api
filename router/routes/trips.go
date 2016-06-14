package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func TripsHandler(w http.ResponseWriter, r *http.Request) {
  trips := []models.Trip{}
  models.DB.Find(&trips)
  json.NewEncoder(w).Encode(trips)
}

// TODO: the distance "threshold" filter needs to be changed to a neighborhood filter
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
  if err != nil {fmt.Println(err)}

  user := GetUserFromRequest(r)
  trip.Owner = int(user.ID)
  models.DB.Create(&trip)
  json.NewEncoder(w).Encode(trip)
}

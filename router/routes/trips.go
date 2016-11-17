package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "github.com/jinzhu/gorm"
       "fmt"
       "lightupon-api/googleMaps"
       "lightupon-api/redis"
       )

func TripsHandler(w http.ResponseWriter, r *http.Request) {
  lat, lon := GetUserLocationFromRequest(r)
  trips := []models.Trip{}
  models.DB.Preload("User").Preload("Scenes", func(DB *gorm.DB) *gorm.DB {
    return DB.Order("Scenes.scene_order ASC") // Preload and order scenes for the map view
  }).Order("((latitude - " + lat + ")^2.0 + ((longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc;").Find(&trips)

  // TODO: this should probably be abstracted in some way
  // Now attach locations from redis
  for i, trip := range trips {
    // Ok what we're going to do here is attempt to pull smooth locations out of Redis. If we're not successful, then attempt to put them in there and try pulling them out again. 
    smoothLocations := redis.GetSmoothedLocationsFromRedis(int(trip.ID))
    if (len(smoothLocations) > 0) {
      trips[i].Locations = smoothLocations
    } else {
      fmt.Println("ERROR: (Attempt 1) Didn't find any smooth locations in Redis for TripID = " + strconv.Itoa(int(trips[i].ID)))

      // see if there are actually any raw locations to be smoothed
      rawLocations := []models.Location{}
      models.DB.Where("trip_id = ?", trips[i].ID).Find(&rawLocations)
      if (len(rawLocations) > 0) {
        smoothLocationsFromGoogle := googleMaps.SmoothTrip(int(trips[i].ID), rawLocations)
        // Now let's throw the smoothed trip up into the DB
        redis.SaveSmoothedLocationsToRedis(int(trips[i].ID), smoothLocationsFromGoogle) //comment this out while testing the GET below
        smoothLocations := redis.GetSmoothedLocationsFromRedis(int(trips[i].ID))
        if (len(smoothLocations) > 0) {
          trips[i].Locations = smoothLocations
        } else {
          fmt.Println("ERROR: (Attempt 2) Didn't find any smooth locations in Redis for TripID = " + strconv.Itoa(int(trips[i].ID)))
        }
      } else {
        fmt.Println("INFO: Didn't find any raw locations in DB for TripID = " + strconv.Itoa(int(trips[i].ID)))
      }
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
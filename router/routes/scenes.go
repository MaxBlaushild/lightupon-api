package routes

import(
       "net/http"
       "lightupon-api/models"
       "lightupon-api/services/redis"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"

       )


// http://localhost:5000/lightupon/scenes/nearby?lat=42.355228&lon=-71.067772
// This endpoint has a little bit of weird behavior. In case A, the user has an active scene and is still at that
// scene. Case B is everything else (the user eith has no active scene or is no longer at that scene).
// In all cases, we will return a scene, along with a flag indicating case A or B. In case B, the scene will be 
// populated with a suggestion for the scene name and nothing else.
func ActiveSceneHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  user.SetUserLocationFromRequest(r)
  scene := user.GetActiveSceneOrSuggestion()
  json.NewEncoder(w).Encode(scene)
}

func PopularScenesHandler(w http.ResponseWriter, r *http.Request) {
  scenes := []models.Scene{}
  models.DB.Where("Featured = true").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

func ScenesIndexHandler(w http.ResponseWriter, r *http.Request) {
  scenes := []models.Scene{}
  models.DB.Preload("Trip.User").Preload("Cards").Order("created_at desc").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

func ScenesHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["tripId"])
  scenes := []models.Scene{}
  models.DB.Find(models.Trip{}, id).Association("Scenes").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

func CreateSelfieSceneHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeTrip := user.ActiveTrip()

  selfie := models.Selfie{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&selfie)

  currentScene := getSceneFromCache(activeTrip.ID)
  currentLocation := models.UserLocation{Latitude: selfie.Location.Latitude, Longitude: selfie.Location.Longitude}
  isAtCurrentScene := currentScene.IsAtScene(currentLocation)

  if (isAtCurrentScene) {
    selfieCard := models.Card{ NibID: "PictureHero", ImageURL: selfie.ImageUrl } 
    currentScene.AppendCard(selfieCard); if err != nil {
      respondWithBadRequest(w, "That selfie was shit!")
      return
    }
  } else {
    selfieScene := models.CreateSelfieScene(selfie)
    err = activeTrip.AppendScene(&selfieScene); if err != nil {
      respondWithBadRequest(w, "That selfie was shit!")
      return
    }
    cacheCurrentScene(selfieScene)
  }

  respondWithCreated(w, "The selfie was created")
}

func getSceneFromCache(tripID uint) (scene models.Scene) {
  key := "currentScene_" + strconv.Itoa(int(tripID))
  redisResponseBytes := redis.GetByteArrayFromRedis(key)
  _ = json.Unmarshal(redisResponseBytes, &scene)
  return
}

func cacheCurrentScene(scene models.Scene) {
  value, _ := json.Marshal(scene)
  key := "currentScene_" + strconv.Itoa(int(scene.TripID))
  redis.SaveByteArrayToRedis(key, value)
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

func ModifySceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneIDint, _ := strconv.Atoi(vars["sceneID"])
  sceneID := uint(sceneIDint)
  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&scene)
  if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk.")
  }
  scene.ID = sceneID

  // TODO iterate through fields instead of doing this one-by-one
  if (scene.Name != "") {models.DB.Model(&scene).Update("name", scene.Name)}
  if (scene.Latitude != 0) {models.DB.Model(&scene).Update("Latitude", scene.Latitude)}
  if (scene.Longitude != 0) {models.DB.Model(&scene).Update("Longitude", scene.Longitude)}
  if (scene.BackgroundUrl != "") {models.DB.Model(&scene).Update("BackgroundUrl", scene.BackgroundUrl)}

  respondWithNoContent(w, "The scene was modified.")
}
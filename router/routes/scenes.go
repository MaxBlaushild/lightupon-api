package routes

import(
       "net/http"
       "lightupon-api/models"
       "lightupon-api/services/redis"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "fmt"
       )

func NearbyScenesHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  lat, lon := GetUserLocationFromRequest(r)
  radius := getStringFromRequest(r, "radius", "0.01")
  numScenes, _ := strconv.Atoi(getStringFromRequest(r, "numScenes", "100"))

  scenes, err := models.GetScenesNearLocation(lat, lon, user.ID, radius, numScenes)
  models.MarkScenesRequest(lat, lon, user.ID, "NearbyScenesHandler")

  if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(scenes)
  }
}

func SceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneID := vars["sceneID"]
  scene, err := models.GetSceneByID(sceneID); if err != nil {
    respondWithBadRequest(w, "ID was bad.")
    return
  }
  json.NewEncoder(w).Encode(scene)
}

func NearbyFollowingScenesHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  lat, lon := GetUserLocationFromRequest(r)
  scenes := models.GetFollowingScenesNearLocation(lat, lon, user.ID)
  models.MarkScenesRequest(lat, lon, user.ID, "NearbyFollowingScenesHandler")
  json.NewEncoder(w).Encode(scenes)
}

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
  user := GetUserFromRequest(r)
  scenes := models.IndexScenes()
  user.SetUserLikenessOfScenes(scenes)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(scenes)
}

func ScenesForUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  userID := vars["userID"]
  scenes := models.GetScenesForUser(userID)
  user.SetUserLikenessOfScenes(scenes)
  json.NewEncoder(w).Encode(scenes)
}

func ScenesHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := strconv.Atoi(vars["tripId"])
  scenes := []models.Scene{}
  models.DB.Find(models.Trip{}, id).Association("Scenes").Find(&scenes)
  json.NewEncoder(w).Encode(scenes)
}

func AppendSceneHandler(w http.ResponseWriter, r *http.Request) {
  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)

  err := decoder.Decode(&scene); if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk!")
    return
  }

  user := GetUserFromRequest(r)
  activeTrip := user.ActiveTrip()

  err = activeTrip.AppendScene(&scene); if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk!")
    return
  }

  json.NewEncoder(w).Encode(scene)
}

func PutSceneHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeTrip := user.ActiveTrip()
  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&scene); if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk!")
    return
  }

  activeTrip.PutScene(&scene)
  cacheCurrentScene(scene)

  respondWithCreated(w, "The scene was created")
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

func CreateSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  tripID := vars["tripID"]
  trip := models.Trip{}
  models.DB.First(&trip, tripID)

  scene := models.Scene{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&scene); if err != nil {
    respondWithBadRequest(w, "That trip ID was not found.")
    return
  }

  err = trip.AppendScene(&scene); if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk!")
    return
  }

  json.NewEncoder(w).Encode(scene)
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
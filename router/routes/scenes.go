package routes

import(
       "net/http"
       "lightupon-api/models"
       "lightupon-api/services/redis"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func NearbyScenesHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  lat, lon := GetUserLocationFromRequest(r)
  radius := getStringFromRequest(r, "radius", "10000")
  numScenes, _ := strconv.Atoi(getStringFromRequest(r, "numScenes", "100"))

  scenes, err := models.GetScenesNearLocation(lat, lon, user.ID, radius, numScenes)
  models.MarkScenesRequest(lat, lon, user.ID, "NearbyScenesHandler")

  if err != nil {
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(scenes)
  }
}

func SceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  sceneID := vars["sceneID"]
  scene, err := models.GetSceneByID(sceneID); if err != nil {
    respondWithBadRequest(w, "ID was bad.")
    return
  }

  scene.SetPercentDiscovered(user.ID)
  json.NewEncoder(w).Encode(scene)
}

func FollowingScenesHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  pageString := getStringFromRequest(r, "page", "0")
  page, _ := strconv.Atoi(pageString)
  scenes := models.GetFollowingScenes(user.ID, page)
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

func FlagSceneHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  sceneID, err1 := GetUIntFromVars(r, "sceneID"); if err1 != nil {
    respondWithBadRequest(w, "ID was bad.")
    return
  }

  decoder := json.NewDecoder(r.Body)
  postParams := struct {Description string}{}
  err2 := decoder.Decode(&postParams)
  if err2 != nil {
    respondWithBadRequest(w, "The flag description you supplied was somehow fucked up.")
  }

  models.DB.Create(&models.Flag{UserID : user.ID, SceneID : sceneID, Description : postParams.Description})
  respondWithCreated(w, "The scene was flagged")
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

func DiscoverSceneHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  sceneID := vars["sceneID"]
  decoder := json.NewDecoder(r.Body)

  err := decoder.Decode(&user.Location); if err != nil {
    respondWithBadRequest(w, "The location sent was bunk.")
    return
  }

  scene, err := models.GetSceneByID(sceneID); if err != nil {
    respondWithBadRequest(w, "The scene you sent us was bunk.")
    return
  }
  
  user.Discover(&scene)
  respondWithNoContent(w, "Explored, my friend.")
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
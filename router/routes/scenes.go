package routes

import(
       "net/http"
       "lightupon-api/models"
       "lightupon-api/services/redis"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       "github.com/kr/pretty"
       )

func NearbyScenesHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  lat, lon := GetUserLocationFromRequest(r)  
  scenes := models.GetScenesNearLocation(lat, lon, user.ID)

  // experimental business. nothing to see here move along..
  // user.UpdateDarknessState(lat, lon) // Update that sweet sweet user state

  models.MarkScenesRequest(lat, lon, user.ID, "NearbyScenesHandler")
  json.NewEncoder(w).Encode(scenes)
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
  pretty.Println(scene)
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
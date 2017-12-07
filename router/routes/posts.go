package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "strconv"
       "fmt"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  post := models.Post{}

  err := decoder.Decode(&post); if err != nil {
    respondWithBadRequest(w, "The post you sent us was bunk.")
    return
  }

  user := GetUserFromRequest(r)

  err = models.DB.Model(&user).Association("Posts").Append(post).Error; if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  }

  json.NewEncoder(w).Encode(post)
}

func GetNearbyPosts(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
  lat, lon := GetUserLocationFromRequest(r)
  radius := getStringFromRequest(r, "radius", "10000")
  numScenes, _ := strconv.Atoi(getStringFromRequest(r, "numScenes", "100"))
  scenes, err := models.GetPostsNearLocation(lat, lon, user.ID, radius, numScenes)

  if err != nil {
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(scenes)
  }
}
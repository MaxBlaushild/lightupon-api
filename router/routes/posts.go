package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
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
  fmt.Println(post.Name)
  err = models.DB.Model(&user).Association("Posts").Append(post).Error; if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  }

  json.NewEncoder(w).Encode(post)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  postID, _ := vars["postID"]
  post, err := models.GetPostByID(postID)

  if err != nil {
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(post)
  }
}

func GetNearbyPosts(w http.ResponseWriter, r *http.Request) {
  lat, lon := GetLocationFromRequest(r)
  radius := getStringFromRequest(r, "radius", "10000")
  numScenes, _ := strconv.Atoi(getStringFromRequest(r, "numScenes", "100"))
  posts, err := models.GetPostsNearLocation(lat, lon, radius, numScenes)

  if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(posts)
  }
}

func GetUsersPosts(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  userID, _ := vars["userID"]
  posts, err := models.GetUsersPosts(userID)

  if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(posts)
  }
}
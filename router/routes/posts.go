package routes

import(
       "lightupon-api/models"
       "net/http"
       "encoding/json"
       "github.com/gorilla/mux"
       // "strconv"
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

func CompletePostHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  postID, _ := vars["postID"]

  err := models.CompletePostForUser(postID, user)

  if err != nil {
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    respondWithAccepted(w, "Post marked as completed.")
  }
}

func GetNearbyPostsAndTryToDiscoverThem(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)

  lat, lon := GetLocationFromRequest(r)
  user := GetUserFromRequest(r)
  radius := GetStringFromRequest(r, "radius", "5000")

  err := decoder.Decode(&user.Location); if err != nil {
    respondWithBadRequest(w, "The location sent was bunk.")
    return
  }

  posts, err := models.GetNearbyPostsAndTryToDiscoverThem(user, lat, lon, radius, 20)

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
package routes

import(
       "lightupon-api/models"
       "net/http"
       "encoding/json"
       "github.com/gorilla/mux"
       // "strconv"
       "fmt"
       "github.com/davecgh/go-spew/spew"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  post := models.Post{}

  post.CreateNewQuestAndSetFieldsOnPost()

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
  lat, lon := GetLocationFromRequest(r)
  user := GetUserFromRequest(r)
  radius := GetStringFromRequest(r, "radius", "5000")

  posts, err := models.GetNearbyPostsAndTryToDiscoverThem(user, lat, lon, radius, 20)

  // For now, let's spew all the scenes out for debugging, to make sure we're getting the right ones.
  spew.Dump(posts)

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
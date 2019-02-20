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

  user := newRequestManager(r).GetUserFromRequest()
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

// TODO: remove this because it's old
func GetNearbyPosts(w http.ResponseWriter, r *http.Request) {
  requestManager := newRequestManager(r)
  lat, lon := requestManager.GetLocationFromRequest()
  user := requestManager.GetUserFromRequest()
  radius := requestManager.getStringFromRequest("radius", "10000")
  numPosts, _ := strconv.Atoi(requestManager.getStringFromRequest("numScenes", "100")) // TODO: clean up numScenes in conjunction with client app

  posts, err := models.GetPostsNearLocationWithUserDiscoveries(lat, lon, user.ID, radius, numPosts)

  if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(posts)
  }
}

func GetNearbyPostsRoute(w http.ResponseWriter, r *http.Request) {
  requestManager := newRequestManager(r)

  posts, err := GetNearbyPostsWithDependencies(requestManager)

  if err != nil {
    fmt.Println(err)
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(posts)
  }
}

func GetNearbyPostsWithDependencies(requestAccessor requestAccessor) (posts []models.Post, err error) {
  lat, lon := requestAccessor.GetLocationFromRequest()
  user := requestAccessor.GetUserFromRequest()
  radius := requestAccessor.getStringFromRequest("radius", "5000")

  posts, err = models.GetPostsNearLocation_NEW(lat, lon, user.ID, radius)

  return
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
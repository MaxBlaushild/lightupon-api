package routes

import(
       "lightupon-api/app"
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

// TODO: remove this because it's old
// func GetNearbyPosts(w http.ResponseWriter, r *http.Request) {
//   lat, lon := GetLocationFromRequest(r)
//   user := GetUserFromRequest(r)
//   radius := GetStringFromRequest(r, "radius", "10000")
//   numPosts, _ := strconv.Atoi(GetStringFromRequest(r, "numScenes", "100")) // TODO: clean up numScenes in conjunction with client app

//   posts, err := models.GetNearbyPosts(lat, lon, user.ID, radius, numPosts)

//   if err != nil {
//     fmt.Println(err)
//     respondWithBadRequest(w, "Something went wrong.")
//   } else {
//     json.NewEncoder(w).Encode(posts)
//   }
// }

func GetNearbyPostsRoute(w http.ResponseWriter, r *http.Request) {
  databaseManager := models.CreateNewDatabaseManager(models.DB)

  // These need to be accessed here because only the router package knows about the requestAccessor
  lat, lon := GetLocationFromRequest(r)
  user := GetUserFromRequest(r)
  radius := GetStringFromRequest(r, "radius", "5000")

  posts, err := app.GetNearbyPosts(user.ID, lat, lon, radius, 20, databaseManager)
  

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
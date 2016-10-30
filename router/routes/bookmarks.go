package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/davecgh/go-spew/spew"
       "fmt"
       "encoding/json"

       )


func LikeBookmark(w http.ResponseWriter, r *http.Request) {
  bookmarkID := GetUIntFromVars(r, "bookmarkID")
  user := GetUserFromRequest(r)
  userID := user.ID
  like := models.Like{BookmarkID:bookmarkID, UserID: userID}
  fmt.Println("like in LikeBookmark")
  spew.Dump(like)
  models.DB.Create(&like)

  




  // decoder := json.NewDecoder(r.Body)
  // err := decoder.Decode(&trip)
  // if err != nil {
  //   respondWithBadRequest(w, "The trip credentials you sent us were wack!")
  // }

  // user := GetUserFromRequest(r)
  // trip.UserID = user.ID
  // bookmark := models.Bookmark{}
  // models.DB.First(&bookmark)
  // json.NewEncoder(w).Encode(bookmark)
  // fmt.Println("ok here's the user we got from the request")
  respondWithNoContent(w, "The bookmark was marked.")
}


func GetBookmarksForUser(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  likes := models.GetLikesForUser(user.ID)

  // Get all the bookmarks that the user should get ('should' to be redefined in the future...)
  bookmarks := []models.Bookmark{}
  models.DB.Find(&bookmarks)
  // Now go through each one and say whether the user has previously liked it
  for i, bookmark := range bookmarks {
    fmt.Printf("ID: %s Age: %d\n", bookmark.ID, bookmark.URL)
    // fmt.Printf("Addr: %p\n", &bookmark)

    // fmt.Println("")
    for _, like := range likes {
      if like.BookmarkID == bookmark.ID {
        fmt.Println("math!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
        bookmarks[i].Liked = true
      }
    }
  }
  fmt.Println("heres the bookmarkssss 9999999999999999999999999999999999999999999999999999999")
  spew.Dump(bookmarks)
  json.NewEncoder(w).Encode(bookmarks)
}

func ServeBookmarks(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/bookmarks.html")
}


func Login(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/login.html")
}
package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/davecgh/go-spew/spew"
       "fmt"
       "encoding/json"
       "strconv"

       )


func LikeBookmark(w http.ResponseWriter, r *http.Request) {
  bookmarkID := GetUIntFromVars(r, "bookmarkID")
  user := GetUserFromRequest(r)
  userID := user.ID
  like := models.Like{BookmarkID:bookmarkID, UserID: userID}
  fmt.Println("like in LikeBookmark")
  spew.Dump(like)
  models.DB.Create(&like)
  respondWithNoContent(w, "The bookmark was marked.")
}

func FuckThisBookmark(w http.ResponseWriter, r *http.Request) {
  bookmarkID := GetUIntFromVars(r, "bookmarkID")
  user := GetUserFromRequest(r)
  userID := user.ID
  fuckThis := models.FuckThis{BookmarkID:bookmarkID, UserID: userID}
  fmt.Println("Fuck this bookmark")
  spew.Dump(bookmarkID)
  spew.Dump(userID)
  spew.Dump(fuckThis)
  models.DB.Create(&fuckThis)
  respondWithNoContent(w, "The bookmark was marked.")
}

func GetBookmarksForUser(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  likes := models.GetLikesForUser(user.ID)
  fuckThises := models.GetFuckThisesForUser(user.ID)

  exclusionList := ""
  for _, fuckThis := range fuckThises {
    exclusionList = exclusionList + strconv.Itoa(int(fuckThis.BookmarkID)) + ","
  }

  // Get all the bookmarks that the user should get ('should' to be redefined in the future...)
  bookmarks := []models.Bookmark{}
  models.DB.Where("id not in (" + exclusionList + "-1)").Order("created_at desc;").Find(&bookmarks) // The -1 is a hack to make the list thing work

  // Now go through each one and say whether the user has previously liked it
  for i, bookmark := range bookmarks {
    for _, like := range likes {
      if like.BookmarkID == bookmark.ID {
        bookmarks[i].Liked = true
      }
    }
  }
  json.NewEncoder(w).Encode(bookmarks)
}

func ServeStatsPage(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/stats.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/login.html")
}
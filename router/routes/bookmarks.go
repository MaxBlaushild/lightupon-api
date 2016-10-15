package routes

import(
       "net/http"
       "lightupon-api/models"
       // "encoding/json"
       "github.com/gorilla/mux"
       // "strconv"
       // "github.com/jinzhu/gorm"
       )


func MarkBookmarkSemiprivate(w http.ResponseWriter, r *http.Request) {
  // request body should consist of {Title: "Balls"}
  vars := mux.Vars(r)
  bookmarkID := vars["bookmarkID"]
  like := models.Like{BookmarkID:bookmarkID, UserID: 1}
  models.DB.Create(&like)
  // trip := models.Trip{}
  // models.DB.First(&trip, id)
  // json.NewEncoder(w).Encode(trip)




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
  respondWithNoContent(w, "The bookmark was marked.")
}

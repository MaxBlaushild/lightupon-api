package app

import(
       // "net/http"
       "lightupon-api/models"
       "github.com/davecgh/go-spew/spew"
       // "encoding/json"
       // "github.com/gorilla/mux"
       // "strconv"
       // "fmt"
)

// TODO: Want to refactor to get rid of this and just pass the accessor to GetPostsNearLocation, but the accessor is defined on the routes package, and that would create a circular dependency
func GetNearbyPosts(lat string, lon string, userID uint, radius string, numPosts int, databaseAccessor models.DatabaseAccessor) (posts []models.Post, err error) {
  firstPosts, err := databaseAccessor.GetFirstPostsNearLocation(lat, lon, radius, numPosts)

  spew.Dump(firstPosts)
  

  // posts, err = models.GetPostsNearLocation_NEW(lat, lon, userID, radius)

  return
}
package app

import(
       // "net/http"
       "lightupon-api/models"
       // "encoding/json"
       // "github.com/gorilla/mux"
       // "strconv"
       // "fmt"
)

// TODO: Want to refactor to get rid of this and just pass the accessor to GetPostsNearLocation, but the accessor is defined on the routes package, and that would create a circular dependency
func GetNearbyPostsWithDependencies(lat string, lon string, userID uint, radius string, ModelsDatabaseAccessor models.DatabaseAccessor) (posts []models.Post, err error) {

  

  // posts, err = models.GetPostsNearLocation_NEW(lat, lon, userID, radius)

  return
}
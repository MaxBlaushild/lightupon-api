package app

import (
  "testing"
  "lightupon-api/models"
  // "github.com/davecgh/go-spew/spew"
  "fmt"
  )

type MockDatabaseAccessor struct {
}

func (mda MockDatabaseAccessor) GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []models.Post, err error) {


  posts = []models.Post{
    models.Post{
      UserID : 1, 
      Location : models.Location{
        Latitude : 71.1234,
        Longitude : -42.1234,
      },
    },
  }
  return
}

func TestGetNearbyPosts(t *testing.T) {
  models.Connect()

  databaseManager := models.CreateNewDatabaseManager(models.DB)
  posts, _ := GetNearbyPosts(1, "42.3459129", "-71.0759857", "5000", 20, databaseManager)
  
  for _, k := range posts {
    fmt.Println(k.ID)
  }
}
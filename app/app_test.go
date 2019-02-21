package app

import (
  "testing"
  "lightupon-api/models"
  "github.com/davecgh/go-spew/spew"
  )

type MockDatabaseAccessor struct {
}

func (mda MockDatabaseAccessor) GetFirstPostsNearLocation(lat string, lon string, radius string, numResults int) (posts []models.Post, err error) {
  // post := models.Post{
  //   UserID : 1, 
  //   Location : models.Location{
  //     Latitude : 71.1234,
  //     Longitude : -42.1234,
  //   },
  // }

  posts = []models.Post{
    models.Post{
      UserID : 1, 
      Location : models.Location{
        Latitude : 71.1234,
        Longitude : -42.1234,
      },
    }
  }
  return
}

func TestGetNearbyPosts(t *testing.T) {
   // t.Errorf("aldkfjhg")
  var mockDatabaseAccessor MockDatabaseAccessor

  firstPosts, err := GetNearbyPosts("asdf", "asdf", 1, "343", 20, mockDatabaseAccessor)
  spew.Dump(firstPosts, err)
}
package models

import(
       "testing"
       "net/http"
       "net/http/httptest"
       )

func TestSomething(t *testing.T) {
	Connect(true) // initialize the global "DB" with testMode = true
	// test stuff here...

  // create a trip

  w := httptest.NewRecorder()
  r := http.Request{}

  LightHandler(r, *w)

  
  selfieCard := Card{ NibID: "PictureHero", ImageURL: "testCard" }
  cards := []Card{selfieCard}

  scene := Scene{ 
    Latitude: 45.1, 
    Longitude: 12.2, 
    SceneOrder: 1, 
    Name: "TestScene",
    BackgroundUrl: "testURL",
  }

  scene.Cards = cards




  userID := uint(rand.Float64()*1000000)
  fmt.Println("userID")
  fmt.Println(userID)
  DB.Create(&scene)
// CreateDegenerateTrip(scene, userID)

// DB.CreateTrip
// 	t.Fail()
// 	t.Error("I'm in a bad mood.")

	// 

	// for a given user, insert an active trip and inserts an active scene
}
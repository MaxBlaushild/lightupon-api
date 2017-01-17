package models

import(
       "testing"
       // "math/rand"
       // "fmt"
       "net/http"
       "net/http/httptest"
       ""


       )


// func TestHeader3D(t *testing.T) {
//     resp := httptest.NewRecorder()

//     uri := "/3D/header/?"
//     path := "/home/test"
//     unlno := "997225821"

//     param := make(url.Values)
//     param["param1"] = []string{path}
//     param["param2"] = []string{unlno}

//     req, err := http.NewRequest("GET", uri+param.Encode(), nil)
//     if err != nil {
//             t.Fatal(err)
//     }

//     http.DefaultServeMux.ServeHTTP(resp, req)
//     if p, err := ioutil.ReadAll(resp.Body); err != nil {
//             t.Fail()
//     } else {
//             if strings.Contains(string(p), "Error") {
//                     t.Errorf("header response shouldn't return error: %s", p)
//             } else if !strings.Contains(string(p), `expected result`) {
//                     t.Errorf("header response doen't match:\n%s", p)
//             }
//     }
// }

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
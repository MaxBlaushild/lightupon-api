package models

import(
    "github.com/jinzhu/gorm"
    "fmt"
    "strconv"
    "net/http"
    "encoding/json"
)

type Location struct {
	gorm.Model
	Latitude float64
	Longitude float64
	TripID uint
}

// So yeah this probably isn't the place to be doing any http stuff, but there's too little of it for me to really care about abstracting it away right now
type MapsResponse struct {
  SnappedPoints []struct {
    PlaceID string
    Location Location
  }
}

func RequestSmoothnessFromGoogle(TripID int, rawLocations []Location) (smoothLocations []Location){
  numberOfLocations := len(rawLocations)
  numberOfChunks := (numberOfLocations / 100) + 1

  fmt.Println("RequestSmoothnessFromGoogle for TripID: " + strconv.Itoa(TripID))
  fmt.Println("  numberOfLocations:" + strconv.Itoa(numberOfLocations))
  fmt.Println("  numberOfChunks:" + strconv.Itoa(numberOfChunks))

  for i := 0; i < numberOfChunks; i++ {
    url := ""
    if (i == (numberOfChunks - 1)) {
      url = BuildSmoothingURL(rawLocations[100*i : numberOfLocations])
    } else {
      url = BuildSmoothingURL(rawLocations[100*i : 100*(i + 1)])
    }

    fmt.Println("  url for index(" + strconv.Itoa(i) + "): " + url)


    response := MapsResponse{}
    getJson(url, &response)
    for _, smoothLocation := range response.SnappedPoints {
      smoothLocation := smoothLocation.Location
      smoothLocation.TripID = uint(TripID)
      smoothLocations = append(smoothLocations, smoothLocation)
    }

    fmt.Println("  number of locations in response for index(" + strconv.Itoa(i) + "): " + strconv.Itoa(len(response.SnappedPoints)))
    fmt.Println("  total number of smooth locations for index(" + strconv.Itoa(i) + "): " + strconv.Itoa(len(smoothLocations)))

  }


  
  return
}

func BuildSmoothingURL(oldLocations []Location) string {
  url := ""
  if (len(oldLocations) > 0) {
    url = url + "https://roads.googleapis.com/v1/snapToRoads?key=AIzaSyBS-y6hKLFKiM5yUWIO0AYR5-lrkCZSvp0&path=" // TODO: Probably want to set this as a var on the package
    url = url + strconv.FormatFloat(oldLocations[0].Latitude, 'f', 6, 64) + "," + strconv.FormatFloat(oldLocations[0].Longitude, 'f', 6, 64)
    for _, oldLocation := range oldLocations {
      url = url + "|" + strconv.FormatFloat(oldLocation.Latitude, 'f', 6, 64) + "," + strconv.FormatFloat(oldLocation.Longitude, 'f', 6, 64)
    }
  }
  fmt.Println("url for google maps smoothing:")
  fmt.Println(url)
  return url
}

func getJson(url string, target interface{}) error {  
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  return json.NewDecoder(r.Body).Decode(target)
}

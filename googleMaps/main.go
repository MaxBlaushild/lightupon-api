package googleMaps

import (
    "log"
    "os"
    "googlemaps.github.io/maps"
    "fmt"
    "lightupon-api/models"
    "strconv"
    "net/http"
    "encoding/json"
    "strings"
)

var (
  googleMaps *maps.Client
)

type MapsResponse struct {
  SnappedPoints []struct {
    PlaceID string
    Location models.Location
  }
}

func Init() {
    googleApiKey := os.Getenv("GOOGLE_MAPS_API")
    var err error
    googleMaps, err = maps.NewClient(maps.WithAPIKey(googleApiKey))  // Deprecated but let's let it hang around a little longer

    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }
}

func SmoothTrip(TripID int, rawLocations []models.Location) (smoothLocations []models.Location){
  
  // Build the URL 
  url := BuildSmoothingURL(rawLocations)

  // Build a respond object and reflect onto it
  response := MapsResponse{}
  getJson(url, &response)

  // Now we need to take that response and put all the locations up in the DB
  if (len(response.SnappedPoints) > 0) {
    for _, smoothLocation := range response.SnappedPoints {
      locaish := smoothLocation.Location
      locaish.TripID = uint(TripID)
      smoothLocations = append(smoothLocations, locaish)
    }
  } else {
    fmt.Println("ERROR: No smoothed locations returned from Google for TripID = " + strconv.Itoa(TripID))
  }
  return
}

func BuildSmoothingURL(oldLocations []models.Location) string {
  url := ""
  if (len(oldLocations) > 0) {
    apiKey := strings.Trim(os.Getenv("GOOGLE_MAPS_API"), "⏎")
    url = url + "https://roads.googleapis.com/v1/snapToRoads?key=" + apiKey + "&path=" // TODO: Probably want to set this as a var on the package
    url = url + strconv.FormatFloat(oldLocations[0].Latitude, 'f', 6, 64) + "," + strconv.FormatFloat(oldLocations[0].Longitude, 'f', 6, 64)
    for _, oldLocation := range oldLocations {
      url = url + "|" + strconv.FormatFloat(oldLocation.Latitude, 'f', 6, 64) + "," + strconv.FormatFloat(oldLocation.Longitude, 'f', 6, 64)
    }
  }
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

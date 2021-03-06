package models

import(
    "net/http"
    "encoding/json"
)

var chunkSize = 99

type Location struct {
  Latitude float64
  Longitude float64
  Course float64
  Accuracy float64
  Floor int
  Context string
}

func LocationsAreWithinThreshold(firstLocation Location, secondLocation Location, threshold float64) (isWithinThreshold bool) {
  distance := CalculateDistance(firstLocation, secondLocation)
  isWithinThreshold = distance < threshold
  return
}

func getJson(url string, target interface{}) error {  
  r, err := http.Get(url)
  if err != nil {
    return err
  }
  defer r.Body.Close()

  return json.NewDecoder(r.Body).Decode(target)
}

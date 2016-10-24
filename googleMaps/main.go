package googleMaps

import (
    "log"
    "fmt"
    "os"
    "googlemaps.github.io/maps"
    "golang.org/x/net/context"
    "lightupon-api/models"
)

var (
  googleMaps *maps.Client
)

func Init() {
    googleApiKey := os.Getenv("GOOGLE_MAPS_API")
    var err error
    googleMaps, err = maps.NewClient(maps.WithAPIKey(googleApiKey))

    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }
}

func SnapLocations(locations []models.Location)(newLocations []models.Location) {
    path := locationsToLatLngs(locations)

    request := &maps.SnapToRoadRequest{
        Interpolate: true,
        Path: path,
    }

    snapToRoadResponse, err := googleMaps.SnapToRoad(context.Background(), request); if err != nil {
        fmt.Println(err)
    } else {
        newLocations = snappedPointsToLocations(snapToRoadResponse.SnappedPoints)
    }
    return
}

func locationsToLatLngs(locations []models.Location)(latLngs []maps.LatLng) {
  for i := 0; i < len(locations); i++ {
      latLng := maps.LatLng{Lat: locations[i].Latitude, Lng: locations[i].Longitude}
      latLngs = append(latLngs, latLng)
  }
  return
}

func snappedPointsToLocations(snappedPoints []maps.SnappedPoint)(locations []models.Location) {
  for i := 0; i < len(snappedPoints); i++ {
      location := models.Location{Latitude: snappedPoints[i].Location.Lat, Longitude: snappedPoints[i].Location.Lng}
      locations = append(locations, location)
  }
  return
}
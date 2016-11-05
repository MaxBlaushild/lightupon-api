package googleMaps

import (
    "log"
    "os"
    "googlemaps.github.io/maps"
    "golang.org/x/net/context"
    "lightupon-api/models"
    "fmt"
       
       "github.com/davecgh/go-spew/spew"
)

var (
  googleMaps *maps.Client
)

func Init() {
    googleApiKey := os.Getenv("GOOGLE_MAPS_API")
    fmt.Println("googleApiKey")
    spew.Dump(googleApiKey)
    var err error
    googleMaps, err = maps.NewClient(maps.WithAPIKey(googleApiKey))

    if err != nil {
        log.Fatalf("fatal error: %s", err)
    }
}

func SnapLocations(locations []models.Location)([]models.Location) {
    path := locationsToLatLngs(locations)
    // fmt.Println("locations")
    // spew.Dump(locations)
    newLocations := []models.Location{}

    fmt.Println("path")
    spew.Dump(path)

    request := &maps.SnapToRoadRequest{
        // Interpolate: true,
        Path: path,
    }

    fmt.Println("request")
    spew.Dump(request)

    // r := &maps.SnapToRoadRequest{
    //   Path: []&maps.LatLng{
    //     &maps.LatLng{Lat: -35.27801, Lng: 149.12958},
    //     &maps.LatLng{Lat: -35.28032, Lng: 149.12907},
    //     &maps.LatLng{Lat: -35.28099, Lng: 149.12929},
    //     &maps.LatLng{Lat: -35.28144, Lng: 149.12984},
    //     &maps.LatLng{Lat: -35.28194, Lng: 149.13003},
    //     &maps.LatLng{Lat: -35.28282, Lng: 149.12956},
    //     &maps.LatLng{Lat: -35.28302, Lng: 149.12881},
    //     &maps.LatLng{Lat: -35.28473, Lng: 149.12836},
    //   },
    // }

    // fmt.Println("r")
    // spew.Dump(r)


    // snapToRoadResponse, err := googleMaps.SnapToRoad(context.Background(), r); 

    snapToRoadResponse, err := googleMaps.SnapToRoad(context.Background(), request); 
    // if err == nil {
    //     newLocations = snappedPointsToLocations(snapToRoadResponse.SnappedPoints)
    // }
    fmt.Println("err")
    spew.Dump(err)

    fmt.Println("snapToRoadResponse")
    spew.Dump(snapToRoadResponse)

    return newLocations
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
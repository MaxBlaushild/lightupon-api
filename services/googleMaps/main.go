package googleMaps

import(
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"log"
	"os"
)

func newClient() *maps.Client {
	apiKey := os.Getenv("GOOGLE_MAPS_API")
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return c
}

func GetPlaces(latitude float64, longitude float64) []maps.GeocodingResult {
	c := newClient()
	latLng := maps.LatLng{Lat: latitude, Lng: longitude}

	r := &maps.GeocodingRequest{
		LatLng: &latLng,
	}
   
	reverseGeocodeResponse, err := c.ReverseGeocode(context.Background(), r)

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return reverseGeocodeResponse
}

func GetPrettyPlace(latitude float64, longitude float64) (addressMap map[string]string) {
	places := GetPlaces(latitude, longitude)
	mostAccuratePlace := places[0]

	addressMap = make(map[string]string)
	addressMap["FormattedAddress"] = mostAccuratePlace.FormattedAddress
	addressMap["PlaceID"] = mostAccuratePlace.PlaceID

	for _, addressComponent := range mostAccuratePlace.AddressComponents {
		componentType := addressComponent.Types[0]
		addressMap[componentType] = addressComponent.ShortName
	}

  return
}


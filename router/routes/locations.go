package routes

import("net/http"
       "lightupon-api/models"
       "github.com/kr/pretty"
       "encoding/json")

const locationThreshold float64 = 0.05

func AddLocationHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	decoder := json.NewDecoder(r.Body)
	location := models.Location{}

	err := decoder.Decode(&location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}
	pretty.Println("LOCATION*******************************")
	pretty.Println(location)

	facebookId := GetFacebookIdFromRequest(r)
	currentLocation := models.GetCurrentLocationFromRedis(facebookId)
	locationsAreSamish := models.LocationsAreWithinThreshold(currentLocation, location, locationThreshold)

	if (!locationsAreSamish) {
		pretty.Println("SECONDLOCATION*******************************")
		pretty.Println(location)
		errTwo := user.AddLocationToCurrentTrip(location); if errTwo != nil {
			respondWithBadRequest(w, "There was an error adding the location to the user's current trip.")
			return
		}
		models.SaveCurrentLocationToRedis(facebookId, location)
	}

	respondWithCreated(w, "The location was added to the trip.")
}
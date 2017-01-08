package routes

import("net/http"
       "lightupon-api/models"
       "fmt"
       "encoding/json")

const closeThreshold float64 = 0.05
const farThresh float64 = 0.25

func AddLocationHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	decoder := json.NewDecoder(r.Body)
	location := models.Location{}

	err := decoder.Decode(&location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}

	fmt.Println(location)
	facebookId := GetFacebookIdFromRequest(r)
	currentLocation := models.GetCurrentLocationFromRedis(facebookId)
	fmt.Println("maybe doing the thing")
	if (locationShouldSave(location, currentLocation)) {
		fmt.Println("did the thing")
		errTwo := user.AddLocationToCurrentTrip(location); if errTwo != nil {
			respondWithBadRequest(w, "There was an error adding the location to the user's current trip.")
			return
		}
		models.SaveCurrentLocationToRedis(facebookId, location)
	}

	respondWithCreated(w, "The location was added to the trip.")
}

func locationShouldSave(location models.Location, currentLocation models.Location) bool {
	distance := models.CalculateLocationDistance(currentLocation, location)
	locationsAreFarEnough := distance > closeThreshold
	locationsAreCloseEnough := (distance < farThresh || models.Location{} == currentLocation)
	return (locationsAreCloseEnough && locationsAreFarEnough)
}
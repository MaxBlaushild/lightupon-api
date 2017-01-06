package routes

import("net/http"
       "lightupon-api/models"
       "github.com/kr/pretty"
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
	pretty.Println("LOCATION*******************************")
	pretty.Println(location)

	facebookId := GetFacebookIdFromRequest(r)
	currentLocation := models.GetCurrentLocationFromRedis(facebookId)

	if (locationShouldSave(location, currentLocation)) {
		pretty.Println("SAVINGTHISSHIT*******************************")
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
	pretty.Println("DISTANCE*******************************")
	pretty.Println(distance)
	locationsAreFarEnough := distance > closeThreshold
	pretty.Println("LOCATIONSAREFARENOUGH*******************************")
	pretty.Println(locationsAreFarEnough)
	locationsAreCloseEnough := distance < farThresh
	pretty.Println("LOCATIONSARECLOSEENOUGH*******************************")
	pretty.Println(locationsAreCloseEnough)
	return (locationsAreCloseEnough && locationsAreFarEnough)
}
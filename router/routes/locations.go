package routes

import("net/http"
       "lightupon-api/models"
       "encoding/json")

func AddLocationHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	decoder := json.NewDecoder(r.Body)
	location := models.Location{}

	err := decoder.Decode(&location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}

	errTwo := user.AddLocationToCurrentTrip(location); if errTwo != nil {
		respondWithBadRequest(w, "There was an error adding the location to the user's current trip.")
		return
	}

	respondWithCreated(w, "The location was added to the trip.")
}
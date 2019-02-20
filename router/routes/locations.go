package routes

import("net/http"
       "encoding/json")

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	user := newRequestManager(r).GetUserFromRequest()
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user.Location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}

	// NOTE: I'd really like to use dependency injection here in order to create unit tests for the explore function. So it would be user.Explore(databaseAccessor).
	err = user.Explore(); if err != nil {
		respondWithBadRequest(w, "You goofed.")
	}

	respondWithCreated(w, "Did the thing.")
}
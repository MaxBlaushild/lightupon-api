package routes

import(
		"net/http"
		"encoding/json"

       )

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user.Location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}

	// NOTE: I'd really like to use dependency injection here in order to create unit tests for the explore function. So it would be user.Explore(ModelsDatabaseAccessor).
	err = user.TryToDiscoverPosts(); if err != nil {
		respondWithBadRequest(w, "You goofed.")
	}

	respondWithCreated(w, "Did the thing.")
}
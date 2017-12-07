package routes

import("net/http"
       "encoding/json")

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromRequest(r)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user.Location); if err != nil {
		respondWithBadRequest(w, "The location sent was bunk.")
		return
	}

	err = user.Explore(); if err != nil {
		respondWithBadRequest(w, "You goofed.")
	} 

	respondWithCreated(w, "Did the thing.")
}
package routes

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondWithOK(w, "We're healthy!")
}
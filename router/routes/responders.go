package routes

import (
				"net/http"
				)

func respondWithNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(makeMessage(message))
}

func makeMessage(message string) []byte {
	return []byte("{\"message\": \"" + message + "\"}")
}
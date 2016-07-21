package routes

import (
				"net/http"
				)

func makeMessage(message string) []byte {
	return []byte("{\"message\": \"" + message + "\"}")
}

func respondWithNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(makeMessage(message))
}

func respondWithCreated(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusCreated)
	w.Write(makeMessage(message))
}

func respondWithBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(makeMessage(message))
}

func respondWithAccepted(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusAccepted)
	w.Write(makeMessage(message))
}

func respondWithNoContent(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(makeMessage(message))
}
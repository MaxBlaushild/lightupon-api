package routes

import "net/http"

func makeMessage(message string) []byte {
	return []byte("{\"message\": \"" + message + "\"}")
}

func respondWithNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(makeMessage(message))
}

func respondWithOK(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
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

func respondeWithInternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(makeMessage(message))
}

func respondeWithForbidden(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write(makeMessage(message))
}

func respondeWithRecordExists(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusFound)
	w.Write(makeMessage(message))
}
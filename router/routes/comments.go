package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func ScenesCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  sceneIDString := vars["sceneID"]
  sceneID, err := strconv.Atoi(sceneIDString)
    if err != nil {
  	respondWithBadRequest(w, "The trip ID you sent us was bunk.")
  } else {
  	 comments := models.GetCommentsForScene(sceneID)
  	json.NewEncoder(w).Encode(comments)
  }
} 

func TripsCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  tripIDString := vars["tripID"]
  tripID, err := strconv.Atoi(tripIDString)

  if err != nil {
  	respondWithBadRequest(w, "The trip ID you sent us was bunk.")
  } else {
  	 comments := models.GetCommentsForTrip(tripID)
  	json.NewEncoder(w).Encode(comments)
  }

}
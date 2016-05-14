package routes

import(
       "net/http"
       "trip-advisor-backend/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["sceneId"])
  cards := []models.Card{}
  models.DB.Find(&models.Scene{}, id).Association("Cards").Find(&cards)
  json.NewEncoder(w).Encode(cards)
}
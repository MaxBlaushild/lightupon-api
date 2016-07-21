package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       "strconv"
       )

func CardsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["sceneID"])
  cards := []models.Card{}
  models.DB.Find(&models.Scene{}, id).Association("Cards").Find(&cards)
  json.NewEncoder(w).Encode(cards)
}

// request body should look like {"Text":"pickle shoes","ImageURL":"http://d3gqasl9vmjfd8.cloudfront.net/2f5fd585-6dfa-48b0-9bc5-6b03de931469.png","SceneID":1,"CardOrder":2,"NibID":"PictureHero"}
func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  sceneID, _ := strconv.Atoi(vars["sceneID"])

  card := models.Card{}

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&card)
  if err != nil {
    respondWithBadRequest(w, "The card you sent us was wack. Get that weak shit out of here.")
  }

  card.SceneID = uint(sceneID)

  models.DB.Create(&card)
  respondWithCreated(w, "The card was created.")
}
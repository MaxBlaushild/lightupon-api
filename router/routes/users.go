package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "github.com/gorilla/mux"
       // "strconv"
)

func UserLogisterHandler(w http.ResponseWriter, r *http.Request) {
  jsonUser := models.User{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&jsonUser); if err != nil {
  	respondWithBadRequest(w, "The user you sent us was wack af.")
    return
  }

  user := models.User{}
  models.DB.FirstOrCreate(&user, jsonUser)
  models.DB.Model(&user).Updates(jsonUser)
  json.NewEncoder(w).Encode(user.Token)
}

func UserTokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := vars["facebookId"]
  token := models.RefreshTokenByFacebookId(id)
  json.NewEncoder(w).Encode(token)
}

func LightHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  if err := user.Light(); err != nil {
    respondWithBadRequest(w, "There was an error getting user lit. They must be a heavyweight XD.")
    return
  }

  respondWithCreated(w, "User was lit!")
}

func ExtinguishHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  if err := user.Extinguish(); err != nil {
    respondWithBadRequest(w, "There was an error extinguishing user. Call the terminator.")
    return
  }

  respondWithCreated(w, "User was Extinguished!")
}
package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "fmt"
       "github.com/gorilla/mux"
       // "strconv"
)

func UserLogisterHandler(w http.ResponseWriter, r *http.Request) {
  jsonUser := models.User{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&jsonUser); if err != nil {
  	fmt.Println(err)
  }

  user := models.User{}
  models.DB.FirstOrCreate(&user, jsonUser)
  json.NewEncoder(w).Encode(user.Token)
}

func UserTokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := vars["facebookId"]
  token := models.RefreshTokenByFacebookId(id)
  json.NewEncoder(w).Encode(token)
}
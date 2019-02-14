package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "fmt"
       "github.com/gorilla/mux"
)

func UserLogisterHandler(w http.ResponseWriter, r *http.Request) {
  jsonUser := models.User{}
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&jsonUser); if err != nil {
  	respondWithBadRequest(w, "The user you sent us was wack af.")
    return
  }

  if models.UserIsBlackListed(jsonUser.Token) {
    respondeWithForbidden(w, "User blacklisted.")
  } else {
    models.UpsertUser(&jsonUser)
    json.NewEncoder(w).Encode(jsonUser.Token)
  }
}

func TwitterLoginHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  twitterCreds := models.User{}
  decoder := json.NewDecoder(r.Body)

  err := decoder.Decode(&twitterCreds); if err != nil {
    respondWithBadRequest(w, "The user you sent us was no good.")
    return
  }

  err = user.Update(twitterCreds); if err != nil {
    respondWithBadRequest(w, "The user you sent us was no good.")
    return
  }

  json.NewEncoder(w).Encode(user)
}

func InstagramLoginHandler(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query()
  fmt.Println(query)
  respondWithCreated(w, "farts")
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  userID, _ := vars["userID"]
  user := models.GetUserByID(userID)
  json.NewEncoder(w).Encode(user)
}

func UserTokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, _ := vars["facebookId"]
  token := models.RefreshTokenByFacebookId(id)
  json.NewEncoder(w).Encode(token)
}

func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
  query := r.FormValue("full_name")
  users := models.FindUsers(query)
  json.NewEncoder(w).Encode(users)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  json.NewEncoder(w).Encode(user)
}

func AddDeviceToken(w http.ResponseWriter, r *http.Request) {

  user := GetUserFromRequest(r)
  decoder := json.NewDecoder(r.Body)
  
  device := models.Device{}

  if err := decoder.Decode(&device); err != nil {
    respondWithBadRequest(w, "Fuck that shit you sent us.")
    return
  }

  device.UserID = user.ID
  models.DB.FirstOrCreate(&models.Device{}, &device)
  respondWithCreated(w, "Token was inserted!")

}
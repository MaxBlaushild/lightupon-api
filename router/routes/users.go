package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "fmt"
       "github.com/gorilla/mux"
       // "lightupon-api/feature"
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
  currentUser := GetUserFromRequest(r)
  user.PopulateIsFollowing(&currentUser)
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

func LightHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  if err := user.Light(); err != nil {
    respondWithBadRequest(w, "There was an error getting user lit.")
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

func MeHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  user.ManaTotal = models.GetManaTotalForUser(user.ID)
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

func GetManaTotalHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  json.NewEncoder(w).Encode(struct {ManaTotal int}{models.GetManaTotalForUser(user.ID)})
}

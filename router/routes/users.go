package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
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

  models.UpsertUser(&jsonUser)
  json.NewEncoder(w).Encode(jsonUser.Token)
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
  decoder := json.NewDecoder(r.Body)
  location := models.Location{}


  // if (feature.IsFeatureEnabledForUser("drop_stuff_instead_of_selfie", user.ID)) {
  //   models.CreateStuffTrip(user.ID)
  // }

  err := decoder.Decode(&location); if err != nil {
    respondWithBadRequest(w, "The location sent was bunk.")
    return
  }

  if err := user.Light(location); err != nil {
    respondWithBadRequest(w, "There was an error getting user lit. They must be a heavyweight XD.")
    return
  }

  respondWithCreated(w, "User was lit!")
}

func ExtinguishHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  decoder := json.NewDecoder(r.Body)
  location := models.Location{}

  err := decoder.Decode(&location); if err != nil {
    respondWithBadRequest(w, "The location sent was bunk.")
    return
  }

  if err := user.Extinguish(location); err != nil {
    respondWithBadRequest(w, "There was an error extinguishing user. Call the terminator.")
    return
  }

  respondWithCreated(w, "User was Extinguished!")
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  json.NewEncoder(w).Encode(user)
}

func AddDeviceToken(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  vars := mux.Vars(r)
  deviceTokenString, _ := vars["deviceToken"]

  deviceTokenModel := models.DeviceToken{}
  deviceTokenModel.UserID = user.ID
  deviceTokenModel.DeviceToken = deviceTokenString

  models.DB.FirstOrCreate(&deviceTokenModel)
}

package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/context"
       "github.com/dgrijalva/jwt-go"
       "encoding/json"
       "strconv"
       "github.com/gorilla/mux"
       "lightupon-api/websockets"
       "fmt"
       )

func CreatePartyHandler(w http.ResponseWriter, r *http.Request) {  
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip)

  if err != nil {
    respondWithBadRequest(w, "The trip credentials you sent are no bueno!")
  }

  user := GetUserFromRequest(r)
  party := models.Party{TripID: trip.ID}
  party.MoveToNextScene()
  models.DB.Model(&user).Association("Parties").Append(&party)
  websockets.H.AddUserConnectionToParty(user, party)
  json.NewEncoder(w).Encode(party)
}

func GetPartyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["id"])
  party := models.Party{}
  models.DB.First(&party, partyID)
  json.NewEncoder(w).Encode(party)
}

func FinishPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  fmt.Println("HURRRR")
  activeParty := user.ActiveParty()
  fmt.Println("HURRRR")
  activeParty.DropUser(user)
  websockets.H.DeactivateUserFromParty(user, activeParty.Passcode)
  json.NewEncoder(w).Encode(activeParty)
}

func AddUserToPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  passcode, _ := vars["passcode"]
  party := models.Party{}
  models.DB.Where("passcode = ? AND active = true", passcode).First(&party)

  if (party.ID != 0) {
    models.DB.Model(party).Association("Users").Append(&user)
    websockets.H.AddUserConnectionToParty(user, party)
    json.NewEncoder(w).Encode(party)
  } else {
    notFoundMessage := "The requested party does not exist."
    respondWithNotFound(w, notFoundMessage)
  }
}

func MovePartyToNextSceneHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["partyID"])
  party := models.Party{}
  models.DB.Preload("Scene.Cards").First(&party, partyID)
  party.MoveToNextScene()
  websockets.H.Broadcast <- party
  respondWithAccepted(w, "The party was moved to the next scene.")
}

func LeavePartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  activeParty.DropUser(user)
  websockets.H.DeactivateUserFromParty(user, activeParty.Passcode)
  json.NewEncoder(w).Encode(activeParty)
}

func GetUserFromRequest(r *http.Request)(user models.User){
  token := context.Get(r, "user")
  facebookID := token.(*jwt.Token).Claims["facebookId"].(string)
  models.DB.Where("facebook_id = ?", facebookID).First(&user)
  return
}

func GetUserLocationFromRequest(r *http.Request)(lat string, lon string){
  query := r.URL.Query()
  lat = query["lat"][0]
  lon = query["lon"][0]
  return
}

func GetUsersPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  json.NewEncoder(w).Encode(activeParty)
}

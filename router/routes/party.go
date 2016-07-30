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
  models.DB.Model(&user).Association("Parties").Append(&party)
  respondWithCreated(w, "The party was created.")
}

func GetPartyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["id"])
  party := models.Party{}
  models.DB.First(&party, partyID)
  json.NewEncoder(w).Encode(party)
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

func UpdatePartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["partyID"])
  lat, lon := GetUserLocationFromRequest(r)
  pullResponse := models.UpdatePartyStatus(partyID, user.ID, lat, lon)
  json.NewEncoder(w).Encode(pullResponse)
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
  models.DB.Model(user).Association("Parties").Delete(activeParty)
  activeParty.DeactivateIfEmpty()
  websockets.H.DeactivateUserFromParty(user, activeParty.Passcode)
  json.NewEncoder(w).Encode(activeParty)
}

func GetUserFromRequest(r *http.Request)(user models.User){
  token := context.Get(r, "user")
  facebookID := token.(*jwt.Token).Claims["facebookId"].(string)
  models.DB.Where("facebook_id = ?", facebookID).First(&user)
  return
}

func GetUserLocationFromRequest(r *http.Request)(lat float64, lon float64){
  query := r.URL.Query()
  lat, _ = strconv.ParseFloat(query["lat"][0], 64)
  lon, _ = strconv.ParseFloat(query["lon"][0], 64)
  return
}

func GetUsersPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  json.NewEncoder(w).Encode(activeParty)
}

func PartyManagerHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  passcode, _ := vars["passcode"]

  ws, err := websockets.Upgrader.Upgrade(w, r, nil); if err != nil {
    respondWithBadRequest(w, "You done fucked up. Give us a real passcode.")
    return
  }

  c := &websockets.Connection{Send: make(chan models.PullResponse), WS: ws, Passcode: passcode, User: user}

  websockets.H.Register <- c

  go c.ReadPump()
  c.WritePump()
}


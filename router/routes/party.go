package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/context"
       "github.com/dgrijalva/jwt-go"
       "encoding/json"
       "strconv"
       "fmt"
       "github.com/gorilla/mux"
       "lightupon-api/websockets"
       )

func CreatePartyHandler(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip)
  if err != nil {fmt.Println(err)}

  user := GetUserFromRequest(r)
  party := models.Party{TripID: trip.ID}
  models.DB.Model(&user).Association("Parties").Append(&party)
  json.NewEncoder(w).Encode(party)
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
  models.DB.Where("passcode = ? AND active = true", passcode).First(&party).Association("Users").Append(&user)
  // json.NewEncoder(w).Encode(party)
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
}

func LeavePartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["partyID"])
  party := models.Party{}

  models.DB.First(&party, partyID)
  models.DB.Model(user).Association("Parties").Delete(party)
  party.DeactivateIfEmpty()
  websockets.H.DeactivateUser(user, party.Passcode)
}

func GetUserFromRequest(r *http.Request)(user models.User){
  token := context.Get(r, "user")
  id := token.(*jwt.Token).Claims["facebookId"].(string)
  models.DB.Where("facebook_id = ?", id).First(&user)
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
  activeParty := models.Party{}
  parties := []models.Party{}
  models.DB.Model(&user).Association("Parties").Find(&parties)
  for _, party := range parties {
    if party.Active {
      activeParty = party
    }
  }
  json.NewEncoder(w).Encode(activeParty)
}

func PartyManagerHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  user := GetUserFromRequest(r)
  passcode, _ := vars["passcode"]

  ws, err := websockets.Upgrader.Upgrade(w, r, nil); if err != nil {
    fmt.Println(err)
    return
  }

  c := &websockets.Connection{Send: make(chan models.PullResponse), WS: ws, Passcode: passcode, User: user}

  websockets.H.Register <- c

  go c.ReadPump()
  c.WritePump()
}


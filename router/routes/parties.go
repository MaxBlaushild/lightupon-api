package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       "strconv"
       "github.com/gorilla/mux"
       )

func CreatePartyHandler(w http.ResponseWriter, r *http.Request) {  
  decoder := json.NewDecoder(r.Body)
  trip := models.Trip{}
  err := decoder.Decode(&trip)

  if err != nil {
    respondWithBadRequest(w, "The trip credentials you sent are no bueno!")
  }

  user := GetUserFromRequest(r)
  party, _ := user.StartParty(trip.ID)
  json.NewEncoder(w).Encode(party)
}

func GetPartyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["id"])
  party := models.Party{}
  models.DB.First(&party, partyID)
  json.NewEncoder(w).Encode(party)
}

func EndPartyHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  partyID, _ := strconv.Atoi(vars["partyID"])
  party := models.Party{}
  models.DB.First(&party, partyID)
  err := party.Deactivate(); if err != nil {
    respondWithBadRequest(w, "There was an error deactivating party.")
  } else {
    respondWithAccepted(w, "The party was ended.")
  }
}

func AddUserToPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  vars := mux.Vars(r)
  passcode, _ := vars["passcode"]
  party := models.Party{}
  models.DB.Where("passcode = ? AND active = true", passcode).First(&party)

  if (party.ID != 0) {
    party.AddUser(user)
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
  party.Broadcast()
  respondWithAccepted(w, "The party was moved to the next scene.")
}

func CreatePartyInviteHandler(w http.ResponseWriter, r *http.Request) {
  // Get the partyID from the URL and populate a party model
  // vars := mux.Vars(r)
  // partyID, _ := strconv.Atoi(vars["partyID"])
  // party := models.Party{}
  // models.DB.First(&party, partyID)

  // // Get the user
  // user := models.User{}
  // decoder := json.NewDecoder(r.Body)
  // err := decoder.Decode(&user); if err != nil {
  //   respondWithBadRequest(w, "That userID you sent us was all kinds of fucked up.")
  // }

  // // Create the invitation and put that shit up in the DB
  // invite := models.PartyInvite{UserID: user.ID, PartyID:party.ID}
  // models.DB.Create(&invite)

  // // Notify the invitee of the invitation if they have an open connection
  // c :=ets.Hub.Connections[user.FacebookId]
  // if c != nil {
  //   fmt.Println("The invitee has an active connection! Let's send them a pullResponse.")
  //   pullResponse := models.PullResponse{Passcode: party.Passcode, Party: party, Scene: party.Scene, NextScene: party.NextScene(), PartyInvite:invite}
  //   c.Send <- pullResponse
  // } else {
  //   fmt.Println("The invitee doesn't have an active connection!")
  // }

  // respondWithAccepted(w, "The user has been invited.")
}

func LeavePartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  activeParty.DropUser(user)
  json.NewEncoder(w).Encode(activeParty)
}

func GetUsersPartyHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  json.NewEncoder(w).Encode(activeParty)
}

package routes

import (
  "html/template"
  "net/http"
  "lightupon-api/models"
  "encoding/json"
)

func AllQuestsHandler(w http.ResponseWriter, r *http.Request) {
  t := template.Must(template.ParseFiles("html/quests.html"))

  var quests []models.Quest
  models.DB.Order("id asc").Find(&quests)

  data := struct{Quests []models.Quest}{
    Quests: quests,
  }

  t.Execute(w, data)
}

func TrackQuestHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  questID, err := GetUIntFromVars(r, "questID")

  if err != nil {
    respondWithBadRequest(w, "Bad quest ID given")
    return
  }

  err = user.TrackQuest(questID)

  if err != nil {
    respondWithBadRequest(w, "Can't follow that quest.")
  } else {
    respondWithAccepted(w, "Quest successfully tracked.")
  }
}

func UntrackQuestHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  questID, err := GetUIntFromVars(r, "questID")

  if err != nil {
    respondWithBadRequest(w, "Bad quest ID given")
    return
  }

  err = user.UntrackQuest(questID)

  if err != nil {
    respondWithBadRequest(w, "Can't untrack that quest.")
  } else {
    respondWithAccepted(w, "Quest successfully untracked.")
  }
}

func TrackedQuestsHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)

  trackedQuests, err := user.TrackedQuests()

  if err != nil {
    respondWithBadRequest(w, "Unable to retrieve quests.")
  } else {
    json.NewEncoder(w).Encode(trackedQuests)
  }

}

func EditQuestHandler(w http.ResponseWriter, r *http.Request) {
  questID, _ := GetUIntFromVars(r, "questID")
  questYaml := models.GetQuestYaml(questID)

  data := struct{
    QuestID uint
    QuestYaml string
  }{
    QuestID: questID,
    QuestYaml: questYaml,
  }

  t := template.Must(template.ParseFiles("html/editQuest.html"))
  t.Execute(w, data)
}

func GetQuestJsonHandler(w http.ResponseWriter, r *http.Request) {
  questID, err := GetUIntFromVars(r, "questID")

  if err != nil {
    respondWithBadRequest(w, "Something went wrong.")
  } else {
    json.NewEncoder(w).Encode(models.GetQuestForEditing(questID))

  }
}

func GetQuestWithUserContextHandler(w http.ResponseWriter, r *http.Request) {
  questID, err := GetUIntFromVars(r, "questID")
  user := GetUserFromRequest(r)

  if err != nil {
    respondWithBadRequest(w, "Bad quest ID given")
    return
  }

  quest, err := models.GetQuestWithUserContext(questID, user.ID)

  if err != nil {
    respondWithBadRequest(w, "Something went wrong! You're SOL.")
  } else {
    json.NewEncoder(w).Encode(quest)
  }
}

func UpdateQuestHandler(w http.ResponseWriter, r *http.Request) {
  // For now, just assign ownership to the first user we find
  var user models.User
  models.DB.First(&user)

  questID, _ := GetUIntFromVars(r, "questID")

  decoder := json.NewDecoder(r.Body)
  questYamlStruct := struct{QuestYaml string}{}

  err := decoder.Decode(&questYamlStruct); if err != nil {
    respondWithBadRequest(w, "Couldnt pull the yaml out of the request body.")
    return
  }

  err = models.UpdateQuest(questID, questYamlStruct.QuestYaml, user)

  if (err == nil) {
    respondWithAccepted(w, "success")
  } else {
    respondWithBadRequest(w, "Error: Unable to update quest.")
  }
}
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

func ActiveQuestsHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeQuests, err := user.ActiveQuests()

  if err != nil {
    respondWithBadRequest(w, "Unable to retrieve quests.")
  } else {
    json.NewEncoder(w).Encode(activeQuests)
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
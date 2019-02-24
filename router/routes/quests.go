package routes

import (
  "html/template"
  "net/http"
  "github.com/davecgh/go-spew/spew"
  "lightupon-api/models"
  "encoding/json"
)

func AllQuestsHandler(w http.ResponseWriter, r *http.Request) {

  t := template.Must(template.ParseFiles("html/quests.html"))

  var quests []models.Quest
  models.DB.Find(&quests)

  data := struct{Quests []models.Quest}{
    Quests: quests,
  }

  t.Execute(w, data)
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

func UpdateQuestHandler(w http.ResponseWriter, r *http.Request) {

  var user models.User
  models.DB.First(&user)
  spew.Dump(user)

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
    respondWithBadRequest(w, "The type provided in the uri is invalid.")
  }

  // questYaml := models.GetQuestYaml(questID)

  // data := struct{
  //   QuestID uint
  //   QuestYaml string
  // }{
  //   QuestID: questID,
  //   QuestYaml: questYaml,
  // }

  // t := template.Must(template.ParseFiles("html/editQuest.html"))
  // t.Execute(w, data)
}
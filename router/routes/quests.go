package routes

import (
  "html/template"
  "net/http"
  "github.com/davecgh/go-spew/spew"
  "lightupon-api/models"
)

// type TodoPageData struct {
//   Quests    []models.Quest
// }

func AllQuestsHandler(w http.ResponseWriter, r *http.Request) {

  t := template.Must(template.ParseFiles("html/quests.html"))

  var quests []models.Quest
  models.DB.Find(&quests)

  data := struct{Quests []models.Quest}{
    Quests: quests,
  }

  t.Execute(w, data)
}

type QuestForEditing struct {
  ID uint
  Description string
  Posts []struct {
    Latitude string
    Longitude string
    Caption string
  }
}

func EditQuestHandler(w http.ResponseWriter, r *http.Request) {

  t := template.Must(template.ParseFiles("html/editQuest.html"))

  var quest models.Quest
  models.DB.Where("id = 1").First(&quest)

  var posts []models.Post
  models.DB.Where("quest_id = ?", 1).Order("quest_order asc").Find(&posts)

  spew.Dump(posts)

  var questForEditing QuestForEditing
  questForEditing.Description = quest.Description
  questForEditing.ID = quest.ID

  data := struct{Quest QuestForEditing}{
    Quest: questForEditing,
  }

  t.Execute(w, data)
}
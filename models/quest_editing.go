package models

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "errors"
)

type QuestForEditing struct {
  Description string
  Posts []PostForEditing
}

type PostForEditing struct {
  Latitude float64
  Longitude float64
  Caption string
  ImageUrl string
  ID uint
}

func GetQuestYaml(questID uint) (questYaml string) {
  var quest Quest
  DB.Where("id = ?", questID).First(&quest)

  var posts []Post
  DB.Where("quest_id = ?", questID).Order("quest_order asc").Find(&posts)

  var questForEditing QuestForEditing
  questForEditing.Description = quest.Description

  for _, post := range posts {
    postForEditing := PostForEditing{
      ID: post.ID,
      Latitude: post.Latitude,
      Longitude: post.Longitude,
      Caption: post.Caption,
      ImageUrl: post.ImageUrl,
    }

    questForEditing.Posts = append(questForEditing.Posts, postForEditing)
  }

  bytez, err := yaml.Marshal(&questForEditing); if err != nil {
    fmt.Println("ERROR: Fuuuuuck that quest couldnt yaml serialize.", err)
  }

  questYaml = string(bytez)

  return
}

func UpdateQuest(questID uint, questYaml string, user User) (err error) {
  var quest QuestForEditing
  err = yaml.Unmarshal([]byte(questYaml), &quest); if err != nil {
     err = errors.New("Unable to parse quest yaml!")
     return
  }

  var questOrder uint = 0
  for _, postFromClient := range quest.Posts {

    questOrder += 1
    var post Post

    if postFromClient.ID != 0 {
      DB.Where("id = ?", postFromClient.ID).First(&post)
    }

    post.Latitude = postFromClient.Latitude
    post.Longitude = postFromClient.Longitude
    post.Caption = postFromClient.Caption
    post.QuestOrder = questOrder
    post.UserID = user.ID

    if postFromClient.ID != 0 {
      DB.Save(post)
    } else {
      post.QuestID = questID
      DB.Create(&post)
    }
  }

  return
}


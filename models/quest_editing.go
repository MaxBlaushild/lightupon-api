package models

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "errors"
)

type QuestForEditing struct {
  ID uint
  Description string
  TimeToComplete int
  UserID uint
  Posts []PostForEditing
}

type PostForEditing struct {
  QuestID uint
  Latitude float64
  Longitude float64
  Caption string
  ImageUrl string
  ID uint
}

func GetQuestForEditing(questID uint) (questForEditing QuestForEditing) {
  var quest Quest
  DB.Where("id = ?", questID).First(&quest)

  var posts []Post
  DB.Where("quest_id = ?", questID).Order("quest_order asc").Find(&posts)

  questForEditing.ID = quest.ID
  questForEditing.Description = quest.Description
  questForEditing.TimeToComplete = quest.TimeToComplete
  questForEditing.UserID = quest.UserID

  for _, post := range posts {
    postForEditing := PostForEditing{
      ID: post.ID,
      Latitude: post.Latitude,
      Longitude: post.Longitude,
      Caption: post.Caption,
      ImageUrl: post.ImageUrl,
      QuestID: post.QuestID,
    }

    questForEditing.Posts = append(questForEditing.Posts, postForEditing)
  }

  return
}

func GetQuestYaml(questID uint) (questYaml string) {
  questForEditing := GetQuestForEditing(questID)

  bytez, err := yaml.Marshal(&questForEditing); if err != nil {
    fmt.Println("ERROR: Fuuuuuck that quest couldnt yaml serialize.", err)
  }

  questYaml = string(bytez)

  return
}

func UpdateQuest(questID uint, questYaml string, user User) (err error) {
  var questFromClient QuestForEditing
  err = yaml.Unmarshal([]byte(questYaml), &questFromClient); if err != nil {
     err = errors.New("Unable to parse quest yaml!")
     return
  }

  err = updateQuestInDatabase(questFromClient); if err != nil {
     return
  }

  var questOrder uint = 0
  for _, postFromClient := range questFromClient.Posts {

    questOrder += 1
    var post Post

    if postFromClient.ID != 0 {
      DB.Where("id = ?", postFromClient.ID).First(&post)
    }

    post.Latitude = postFromClient.Latitude
    post.Longitude = postFromClient.Longitude
    post.ImageUrl = postFromClient.ImageUrl
    post.Caption = postFromClient.Caption
    post.QuestID = postFromClient.QuestID
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

func updateQuestInDatabase(questFromClient QuestForEditing) (err error) {
  deletePostsThatArentInTheQuestSentByTheClient(questFromClient)

  var quest Quest
  DB.Where("id = ?", questFromClient.ID).First(&quest)
  if quest.ID == 0 {
    err = errors.New("Couldn't update quest in the database because we couldn't find that shit up in the database.")
    return
  }
  quest.Description = questFromClient.Description
  quest.TimeToComplete = calculateEstimatedTimeToComplete(questFromClient)
  
  DB.Save(&quest)
  return
}

func calculateEstimatedTimeToComplete(questFromClient QuestForEditing) (walkingTime int) {
  // Minus one because we want to iterate over pairs of consecutive posts in the quest, and there are n - 1 pairs.
  var totalDistance float64 = 0.0
  for i := 0; i < len(questFromClient.Posts) - 1; i++ {
    totalDistance += CalculateDistance(Location{Latitude: questFromClient.Posts[i].Latitude, Longitude: questFromClient.Posts[i].Longitude}, Location{Latitude: questFromClient.Posts[i + 1].Latitude, Longitude: questFromClient.Posts[i + 1].Longitude})
  }
  return int((totalDistance / 1.4) / 60.0) // 1.4 meters per second is apparently how fast people so here's our stoichiometry to get to estimate time for the quest in minutes
}

func deletePostsThatArentInTheQuestSentByTheClient(questFromClient QuestForEditing) {
  var postsFromDB []Post
  DB.Where("quest_id = ?", questFromClient.ID).Find(&postsFromDB)
  for _, postsFromDB := range postsFromDB {
    if !postIsInQuestFromClient(postsFromDB, questFromClient) {
      DB.Delete(&postsFromDB)
    }
  }
}

func postIsInQuestFromClient(postsFromDB Post, questFromClient QuestForEditing) bool {
  for _, postFromClient := range questFromClient.Posts {
    if postFromClient.ID == postsFromDB.ID {
      return true
    }
  }
  return false
}
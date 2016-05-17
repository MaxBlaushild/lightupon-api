package models

import(
      "math/rand"
      "github.com/jinzhu/gorm"
      "time"
      )

type Party struct {
  gorm.Model
  TripID uint
  Trip Trip
  SceneID uint
  Scene Scene
  CurrentSceneOrderID int `gorm:"default:0"`
  Passcode string
  Active bool `gorm:"default:true"`
  Started bool `gorm:"default:false"`
  Users []User `gorm:"many2many:partyusers;"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func (p *Party) BeforeCreate() {
  p.setPasscode()
}

func (p *Party) setPasscode() {
  rand.Seed(time.Now().UnixNano())
  
  b := make([]byte, 4)
  for i := range b {
      index := rand.Intn(len(letterBytes))
      b[i] = letterBytes[index]
  }
  p.Passcode = string(b)
}

func (p *Party) NextScene()(nextScene Scene) {
  DB.Preload("Cards").Where("trip_id = ? AND scene_order = ?", p.TripID, p.CurrentSceneID + 1).First(&nextScene)
  return
}

func (p *Party) MoveToNextScene() {
  DB.Model(&p).Update(map[string]interface{}{
    "current_scene_id": p.CurrentSceneOrderID + 1,
    "scene": p.NextScene(),
  })
}

func UpdatePartyStatus(partyID int, userID uint, user_lat float64, user_lon float64)(pullResponse PullResponse){
  partyuser := Partyuser{}
  DB.Preload("Party").Where("user_id = ? AND party_id = ?", userID, partyID).First(&partyuser)
  arrivedAtNextScene, nextScene := partyuser.IsUserAtNextScene(user_lat, user_lon)

  if (arrivedAtNextScene) {
    DB.Model(&partyuser).Update("current_scene_id", nextScene.ID) // Update the current_scene for partyUser

    DB.Where("scene_id = ?", nextScene.ID).Find(&nextScene.Cards) // Poopulate cards for scene

    // Get the current scene for each user, dedupe the list, and check if everyone is on the same scene
    allCurrentScenes := []int64{}
    DB.Model(&Partyuser{}).Where("party_id = ?", partyuser.Party.ID).Pluck("current_scene_id", &allCurrentScenes)
    uniqueCurrentScenes := removeDuplicates(allCurrentScenes)
    if (len(uniqueCurrentScenes) == 1) {
      DB.Model(&partyuser.Party).Update("current_scene_id", uniqueCurrentScenes[0])
    }

    pullResponse.NextScene = nextScene

  }

  pullResponse.NextSceneAvailable = arrivedAtNextScene

  return
}

func (party *Party) DeactivateIfEmpty() {
  users := []User{}
  DB.Model(&party).Association("Users").Find(&users); if len(users) == 0 {
    DB.Model(&party).Update("active", false)
  }
}

func removeDuplicates(elements []int64) []int64 {
    // Use map to record duplicates as we find them.
    encountered := map[int64]bool{}
    result := []int64{}

    for v := range elements {
      if encountered[elements[v]] == true {
          // Do not add duplicate.
      } else {
          // Record this element as an encountered element.
          encountered[elements[v]] = true
          // Append to result slice.
          result = append(result, elements[v])
      }
    }
    // Return the new slice.
    return result
}
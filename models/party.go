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

func (p *Party) DropUser(user User) {
  DB.Model(user).Association("Parties").Delete(p)
  p.DeactivateIfEmpty()
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
  DB.Preload("Cards").Where("trip_id = ? AND scene_order = ?", p.TripID, p.CurrentSceneOrderID + 1).First(&nextScene)
  return
}

func (p *Party) MoveToNextScene() {
  DB.Model(&p).Update(map[string]interface{}{
    "current_scene_order_id": p.CurrentSceneOrderID + 1,
    "scene": p.NextScene(),
  })
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
package models

import(
      "math/rand"
      "github.com/jinzhu/gorm"
      "time"
      "lightupon-api/live"
      )

type Party struct {
  gorm.Model
  TripID uint `gorm:"index"`
  Trip Trip
  SceneID uint `gorm:"index"`
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

func (p *Party) AfterCreate() {
  p.SyncWithLive()
}

func (p *Party) LoadCurrentScene() {
  DB.Where("trip_id = ? AND scene_order = ?", p.TripID, p.CurrentSceneOrderID+1).Find(&p.Scene)
}

func (p *Party) LoadTrip() {
  DB.Model(&p).Association("Trip").Find(&p.Trip)
}

func (p *Party) LiveParty() (liveParty live.Party) {
  scenes := []Scene{}
  DB.Where("trip_id = ?", p.TripID).Order("scene_order asc").Find(&scenes)

  liveParty = live.Party{
    Users: make(map[uint]*live.Connection),
    Passcode: p.Passcode,
    Objectives: scenesToObjectives(scenes),
    CurrentObjectiveIndex: p.CurrentSceneOrderID,
  }
  return
}

func (p *Party) SyncWithLive() {
  liveParty := p.LiveParty()
  live.Hub.PutParty <- liveParty
}

func scenesToObjectives(scenes []Scene) []live.Objective {
    objectives := make([]live.Objective, len(scenes))
    for index, scene := range scenes {
        location := live.Location{Latitude: scene.Latitude, Longitude: scene.Longitude}
        objectives[index] = live.Objective{Location: location}
    }
    return objectives
}

func (p *Party) DropUser(user User) {
  DB.Model(user).Association("Parties").Delete(p)
  live.Hub.DropUserFromParty(user.ID, p.Passcode)
  p.DeactivateIfEmpty()
}

func (p *Party) Connect(c *live.Connection) {
  live.Hub.Register <- c
  go c.ReadPump()
  c.WritePump()
}

func (p *Party) AddUser(user User) (err error) {
  err = DB.Model(p).Association("Users").Append(&user).Error
  live.Hub.AddUserToParty(user.ID, p.Passcode)
  return
}

func (p *Party) Broadcast() {
  live.Hub.Broadcast <- live.Response{ Passcode: p.Passcode }
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
  p.SyncWithLive()
}

func (p *Party) Deactivate() (err error) {
  live.Hub.UnregisterParty(p.Passcode)
  err = DB.Model(&p).Update("active", false).Error
  err = DB.Model(&p).Association("Users").Clear().Error
  return
}

func (party *Party) DeactivateIfEmpty() {
  users := []User{}
  DB.Model(&party).Association("Users").Find(&users); if len(users) == 0 {
    DB.Model(&party).Update("active", false)
    live.Hub.UnregisterParty(party.Passcode)
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

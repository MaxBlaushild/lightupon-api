package models

import(
      "github.com/jinzhu/gorm"
      "lightupon-api/live"
      )

type ExposedScene struct {
  gorm.Model
  UserID uint
  SceneID uint
  Blur float64
  Unlocked bool
}

const unlockThresholdSmall = 0.06
const unlockThresholdLarge = 0.4


func (exposedScene *ExposedScene) upsertExposedScene(newBlur float64, sceneID uint, userID uint, hidden bool) {
  if (exposedScene.ID == 0) {
    DB.Create(&ExposedScene{UserID : userID, SceneID : sceneID, Blur : newBlur, Unlocked : !hidden})
  } else {
    DB.Model(&exposedScene).Update("Blur", newBlur).Update("Unlocked", !hidden)
  }

  sceneUpdate := live.SceneUpdate{UpdatedSceneID: exposedScene.SceneID, UserID: exposedScene.UserID}
  live.Hub.UpdateClient <- sceneUpdate
}

type SceneUpdate struct {
  UpdatedSceneID uint
  UserID uint
}
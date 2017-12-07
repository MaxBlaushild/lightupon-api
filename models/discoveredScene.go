package models

import(
      "github.com/jinzhu/gorm"
      // "lightupon-api/live"
      )

type DiscoveredScene struct {
  gorm.Model
  UserID uint
  SceneID uint
  PercentDiscovered float64
}

// const unlockThresholdSmall float64 = 20
// const unlockThresholdLarge float64 = 200
// const fadePeriod float64 = 4

// func (dS *DiscoveredScene) NotFullyDiscovered() bool {
//   return dS.PercentDiscovered < 1.0
// }

// func UpsertDiscoveredScene(discoveredScene *DiscoveredScene) {
//   if DB.NewRecord(discoveredScene) {
//     DB.Create(&discoveredScene)
//   } else {
//     DB.Model(&discoveredScene).Update("PercentDiscovered", discoveredScene.PercentDiscovered)
//   }
//   sceneUpdate := live.SceneUpdate{UpdatedSceneID: discoveredScene.SceneID, UserID: discoveredScene.UserID}
//   live.Hub.UpdateClient <- sceneUpdate
// }

// func (dS *DiscoveredScene) UpdatePercentDiscovered(user *User, scene *Scene) {
//   newPercentDiscovered := calculatePercentDiscovered(user, scene)
//   if (newPercentDiscovered > dS.PercentDiscovered) {
//     dS.PercentDiscovered = newPercentDiscovered
//     UpsertDiscoveredScene(dS)
//   }
// }

// func calculatePercentDiscovered(user *User, scene *Scene) (percentDiscovered float64) {
//   distance := CalculateDistance(user.Location, UserLocation{Latitude: scene.Latitude, Longitude: scene.Longitude})
//   if (distance < unlockThresholdSmall) {
//     percentDiscovered = 1.0
//   } else if (distance > unlockThresholdLarge) {
//     percentDiscovered = 0.0
//   } else {
//     percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall)) // TODO: Update this to be a nice smoove cosine function
//   }
//   return
// }

// func GetCurrentDiscoveredScene(userID uint, sceneID uint) DiscoveredScene {
//   discoveredScene := DiscoveredScene{UserID: userID, SceneID: sceneID}
//   DB.First(&discoveredScene, discoveredScene)
//   return discoveredScene
// }
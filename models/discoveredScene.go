package models

import(
      "github.com/jinzhu/gorm"
      "lightupon-api/live"
      )

type DiscoveredScene struct {
  gorm.Model
  UserID uint
  SceneID uint
  PercentDiscovered float64
}

const unlockThresholdSmall float64 = 20
const unlockThresholdLarge float64 = 200

func (dS *DiscoveredScene) NotFullyDiscovered() bool {
  return dS.PercentDiscovered < 1.0
}

func UpsertDiscoveredScene(discoveredScene *DiscoveredScene) {
  if DB.NewRecord(discoveredScene) {
    DB.Create(&discoveredScene)
  } else {
    DB.Model(&discoveredScene).Update("PercentDiscovered", discoveredScene.PercentDiscovered)
  }
  sceneUpdate := live.SceneUpdate{UpdatedSceneID: discoveredScene.SceneID, UserID: discoveredScene.UserID}
  live.Hub.UpdateClient <- sceneUpdate
}

func (dS *DiscoveredScene) UpdatePercentDiscovered(user *User, scene *Scene) {
  newPercentDiscovered := calculatePercentDiscovered(user, scene)
  if (newPercentDiscovered > dS.PercentDiscovered) {
    dS.PercentDiscovered = newPercentDiscovered
    UpsertDiscoveredScene(dS)
  }
}

func calculatePercentDiscovered(user *User, scene *Scene) (percentDiscovered float64) {
  distance := CalculateDistance(user.Location, UserLocation{Latitude: scene.Latitude, Longitude: scene.Longitude})
  if (distance < unlockThresholdSmall) {
    percentDiscovered = 1.0
  } else if (distance > unlockThresholdLarge) {
    percentDiscovered = 0.0
  } else {
    // TODO: Update this to be a nice smoove cosine function
    percentDiscovered = 1.0 - ((distance - unlockThresholdSmall) / (unlockThresholdLarge - unlockThresholdSmall))
  }
  return
}

func GetCurrentDiscoveredScene(userID uint, sceneID uint) DiscoveredScene {
  discoveredScene := DiscoveredScene{UserID: userID, SceneID: sceneID}
  DB.First(&discoveredScene, discoveredScene)
  return discoveredScene
}

// func possiblyRecomputeAllDiscovery(lat string, lon string, userID uint) {
//   exposedScene := ExposedScene{}
//   DB.First(&exposedScene)
//   if exposedScene.ID == 0 {
//     recomputeAllDiscovery()
//   }
// }

// func recomputeAllDiscovery() {
//   fmt.Println("NOTICE: Recomputing all discovery!")
//   locations := []Location{}
//   DB.Find(&locations)
//   for i := 0; i < len(locations); i++ {
//     lat := strconv.FormatFloat(locations[i].Latitude, 'E', -1, 64)
//     lon := strconv.FormatFloat(locations[i].Longitude, 'E', -1, 64)
//     scenes := []Scene{}
//     DB.Find(&scenes)
//     for i := 0; i < len(scenes); i++ {
//       scenes[i].discover(locations[i].UserID, lat, lon)
//     }
//   }
// }


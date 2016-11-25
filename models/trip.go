package models

import(
      "github.com/jinzhu/gorm"
      "lightupon-api/redis"
      "strconv"
       "encoding/json"
       "fmt"
      )

type Trip struct {
  gorm.Model
  Title string `gorm:"not null"`
  Description string `gorm:"not null"`
  Details string
  ImageUrl string `gorm:"not null"`
  Distance float32
  EstimatedTime int
  UserID uint
  User User
  Scenes []Scene
  Locations []Location
  Active bool `gorm:"default:true"`
}

func (t *Trip) PutLocations(locations []Location) {
  DB.Model(&t).Association("Locations").Replace(locations)
}

func GetTripsNearLocation(lat string, lon string) (trips []Trip) {
  DB.Preload("User").Preload("Scenes", func(DB *gorm.DB) *gorm.DB {
    return DB.Order("Scenes.scene_order ASC") // Preload and order scenes for the map view
  }).Order("((latitude - " + lat + ")^2.0 + ((longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc;").Find(&trips)

  for i, _ := range trips {
    locations := GetLocationsForTrip(trips[i].ID)
    trips[i].Locations = locations
  }

  return
}

func GetLocationsForTrip(tripID uint) (locations []Location){
  if (redis.GetRedisKey("smoothing_disabled") == "true") {
    fmt.Println("INFO: Smoothing disabled. Incidentally, here's the TripID: " + strconv.Itoa(int(tripID)))
    DB.Where("trip_id = ?", tripID).Find(&locations)
    return
  } else {
    locations = GetSmoothedLocationsFromRedis(int(tripID))

    // If we find something in redis then return
    if (len(locations) > 0) {
      fmt.Println("INFO: Found some smooth locations in Redis for TripID = " + strconv.Itoa(int(tripID)) + ". Total points found was: " + strconv.Itoa(len(locations)))
      return
    } else {
      fmt.Println("INFO: Didn't find any smooth locations in Redis for TripID = " + strconv.Itoa(int(tripID)))
    }

    // Try to pull the locations out of the DB. If we find nothing then we're SOL, so return. If we found something, we might need it later if we fail to smooth
    rawLocations := []Location{}; DB.Where("trip_id = ?", tripID).Find(&rawLocations)  // If we decide later that we never want to display raw trips, then we should just reflect onto 'locations' here
    if (len(rawLocations) == 0) {
      fmt.Println("INFO: Didn't find any raw locations in DB for TripID = " + strconv.Itoa(int(tripID)))
      return
    }

    if (!AllowSmoothingRequestForTrip(tripID)) { 
      fmt.Println("INFO: Smoothing request rate limited for TripID = " + strconv.Itoa(int(tripID)))
      return
    }

    locations = RequestSmoothnessFromGoogle(int(tripID), rawLocations)
    redis.SetRedisKey("smoothing_request_rate_limit_tripID_" + strconv.Itoa(int(tripID)), "x", 86400) // Rate limit to one day 86400

    if (len(locations) == 0) { 
      fmt.Println("ERROR: Didn't get any smooth locations back from Google for TripID = " + strconv.Itoa(int(tripID)))
      return rawLocations // ok if we've tried all that stuff and nothing has worked, just return the raw locations
    } else {
      // AHA! we got some smoothness back from google, save that shit in redis and also return it
      fmt.Println("INFO: We got some smooth locations back from Google for TripID = " + strconv.Itoa(int(tripID)))
      SaveSmoothedLocationsToRedis(tripID, locations)
      return
    }
  }
}

func SaveSmoothedLocationsToRedis(tripID uint, locations []Location) {
  value, _ := json.Marshal(locations) 
  key := "locations_" + strconv.Itoa(int(tripID))
  redis.SaveByteArrayToRedis(key, value) //comment this out while testing the GET below
}

func AllowSmoothingRequestForTrip(tripID uint) bool {
  rate_limit := redis.GetRedisKey("smoothing_request_rate_limit_tripID_" + strconv.Itoa(int(tripID)))
  return !(rate_limit == "x") // If we find anything in Redis, it won't be an empty string, so this will return false
}

func GetSmoothedLocationsFromRedis(TripID int) (smoothLocations []Location) {
  key := "locations_" + strconv.Itoa(TripID)
  redisResponseBytes := redis.GetByteArrayFromRedis(key)
  _ = json.Unmarshal(redisResponseBytes, &smoothLocations)
  return
}


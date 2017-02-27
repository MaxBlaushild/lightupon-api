package models

import(
      "github.com/jinzhu/gorm"
      "lightupon-api/services/redis"
      "strconv"
       "encoding/json"
       "fmt"
       "hash/fnv"
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
  Comments []Comment
  Active bool
  Constellation []ConstellationPoint
  UserHasLikedTrip bool `sql:"-"`
  TotalLikes int `sql:"-"`
}

type ConstellationPoint struct {
  DeltaY float64
  DistanceToPreviousPoint float64
}

func (t *Trip) AppendScene(scene *Scene) (err error) {
  scene.SceneOrder = uint(len(t.Scenes) + 1)
  err = DB.Model(&t).Association("Scenes").Append(scene).Error
  return
}

func (t *Trip) PutScene(scene *Scene) {
  if scene.ID != 0 {
    card := scene.Cards[0]
    scene.Cards = nil
    DB.Model(&scene).Update(scene)
    scene.AppendCard(card)
  } else {
    scene.UserID = t.UserID
    t.AppendScene(scene)
  }
}

func GetTrip(tripID int, userID uint) (trip Trip) {
  DB.Preload("User").Preload("Scenes.Cards").First(&trip, tripID)
  trip.LoadConstellation()
  trip.LoadCommentsForTrip()
  trip.LoadLikeStuff(userID)
  return
}

func GetTripsForUser(userID string) (trips []Trip) {
  DB.Preload("Scenes.Cards").Order("created_at desc").Where("user_id = ?", userID).Find(&trips)
  for i, _ := range trips {
    trips[i].SetLocations()
  }
  return
}

func (t *Trip) PutLocations(locations []Location) {
  DB.Model(&t).Association("Locations").Replace(locations)
}

func GetTripsNearLocation(lat string, lon string, userID uint) (trips []Trip) {

  // DB.Preload("User").Preload("Scenes.Cards").Order("((latitude - " + lat + ")^2.0 + ((longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc;").Find(&trips)
  DB.Preload("User").Preload("Scenes.Cards").Find(&trips)

  for i, _ := range trips {
    trips[i].SetLocations()

    // ok now take the those locations, try to make a constellation out of them, and attach that to the trip
    trips[i].LoadConstellation()
    trips[i].LoadCommentsForTrip()
    trips[i].LoadLikeStuff(userID)
  }

  return
}

// This function accepts a string and returns a float64 between 0 and 1
func hashStringToDecimal(s string) float64 {
  h := fnv.New32()
  h.Write([]byte("xxxxx" + s + "xxxxx"))
  intHash := h.Sum32()

  blurg := float64(intHash) / 100000.0
  return blurg - float64(int(blurg))
}

func (trip *Trip) LoadConstellation() {
  if (trip.Scenes == nil) {
    // Should probably log an error here
    return
  }

  var constellationPoints []ConstellationPoint

  for i, _ := range trip.Scenes {
    constellationPoint :=  ConstellationPoint{}
    if (i != 0) {  // doesn't make sense to do this for the first scene // REMEMBER INDICIES START AT ZERO!!!!!!!
      // CalculateDistance only accepts a UserLocation as opposed to a Location, so that's what we're gonna use
      // TODO: Refactor the entire app to either use UserLocation or Location OR have UserLocation extend Location
      location1, location2 := UserLocation{}, UserLocation{}
      location1.Longitude = trip.Scenes[i - 1].Longitude
      location1.Latitude = trip.Scenes[i - 1].Latitude
      location2.Longitude = trip.Scenes[i].Longitude
      location2.Latitude = trip.Scenes[i].Latitude
      
      constellationPoint.DistanceToPreviousPoint = CalculateDistance(location1, location2)
    }

    constellationPoint.DeltaY = hashStringToDecimal(string(trip.Scenes[i].ID))
    constellationPoints = append(constellationPoints, constellationPoint)
    trip.Scenes[i].ConstellationPoint = constellationPoint
  }

  trip.Constellation = constellationPoints

  return
}

func (trip *Trip) SetLocations() {

  locations := []Location{}

  // Don't do no smoothing if th e trip is ongoing or if the feature is toggled off
  if ((redis.GetRedisKey("smoothing_disabled") == "true") || trip.Active) {
    fmt.Println("INFO: Smoothing disabled. Incidentally, here's the TripID: " + strconv.Itoa(int(trip.ID)) + ". Also this is the active flag for the trip " + strconv.FormatBool(trip.Active))
    DB.Where("trip_id = ?", trip.ID).Find(&locations)
    trip.Locations = locations
    return
  } else {
    locations = GetSmoothedLocationsFromRedis(int(trip.ID))

    // If we find something in redis then return
    if (len(locations) > 0) {
      fmt.Println("INFO: Found some smooth locations in Redis for TripID = " + strconv.Itoa(int(trip.ID)) + ". Total points found was: " + strconv.Itoa(len(locations)))
      return
    } else {
      fmt.Println("INFO: Didn't find any smooth locations in Redis for TripID = " + strconv.Itoa(int(trip.ID)))
    }

    // Try to pull the locations out of the DB. If we find nothing then we're SOL, so return. If we found something, we might need it later if we fail to smooth
    rawLocations := []Location{}; DB.Where("trip_id = ?", trip.ID).Find(&rawLocations)  // If we decide later that we never want to display raw trips, then we should just reflect onto 'locations' here
    if (len(rawLocations) == 0) {
      fmt.Println("INFO: Didn't find any raw locations in DB for TripID = " + strconv.Itoa(int(trip.ID)))
      trip.Locations = locations
      return
    }

    if (!AllowSmoothingRequestForTrip(trip.ID)) { 
      fmt.Println("INFO: Smoothing request rate limited for TripID = " + strconv.Itoa(int(trip.ID)))
      trip.Locations = locations
      return
    }

    locations = RequestSmoothnessFromGoogle(int(trip.ID), rawLocations)
    redis.SetRedisKey("smoothing_request_rate_limit_tripID_" + strconv.Itoa(int(trip.ID)), "x", 86400) // Rate limit to one day 86400

    if (len(locations) == 0) { 
      fmt.Println("ERROR: Didn't get any smooth locations back from Google for TripID = " + strconv.Itoa(int(trip.ID)))
      trip.Locations = rawLocations
      return  // ok if we've tried all that stuff and nothing has worked, just return the raw locations
    } else {
      // AHA! we got some smoothness back from google, save that shit in redis and also return it
      fmt.Println("INFO: We got some smooth locations back from Google for TripID = " + strconv.Itoa(int(trip.ID)))
      SaveSmoothedLocationsToRedis(trip.ID, locations)
      trip.Locations = locations
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

func GetBookmarkCards() []Card {
  cards := []Card{}
  bookmarks := []Bookmark{}
  DB.Limit(5).Order("created_at desc").Find(&bookmarks)
  for i, bookmark := range bookmarks {
    bookmarkCard := Card{ 
      Caption: bookmark.URL,
      CardOrder: uint(i),
      NibID: "TextHero",
    }
    cards = append (cards, bookmarkCard)
  }
  return cards
}

func (trip *Trip) LoadLikeStuff(userID uint) {
  // Find out whether the user has liked this trip
  trip.UserHasLikedTrip = HasUserLikedTrip(userID, trip.ID)

  // Find out the total number of likes for the trip
  trip.TotalLikes = GetTotalLikesForTrip(trip.ID)
}


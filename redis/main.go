package redis

import(
       "lightupon-api/models"
       "encoding/json"
       "strconv"
       "github.com/garyburd/redigo/redis"
       )


func GetSmoothedLocationsFromRedis(TripID int) []models.Location {

  client, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    key := strconv.Itoa(TripID)
    redisResponse, _ := redis.String(client.Do("GET", key))
    redisResponseBytes := []byte(redisResponse)

    locations := []models.Location{}

    _ = json.Unmarshal(redisResponseBytes, &locations)
    
    return locations
}

func SaveSmoothedLocationsToRedis(TripID int, smoothedLocations []models.Location) {
  key := strconv.Itoa(TripID)
  value, _ := json.Marshal(smoothedLocations)
  client, err := redis.Dial("tcp", ":6379")
  if err != nil {
      panic(err)
  }
  defer client.Close()
  client.Do("SET", key, value)
}


func SetRedisKey(key string, value string) {
    client, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer client.Close()
    client.Do("SET", key, value)
}


func GetRedisKey(key string) string {
    client, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer client.Close()
    value, _ := redis.String(client.Do("GET", key))

    return value
}

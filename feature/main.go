package feature

import(
       "lightupon-api/redis"
       "strconv"
       "fmt"
       )

func IsFeatureEnabled(featureName string) bool {
  redisValue := redis.GetRedisKey("feature:" + featureName)
  return (redisValue != "x")
}

func IsFeatureEnabledForUser(featureName string, userID uint) bool {
  redisValue := redis.GetRedisKey("feature:" + featureName + ",user:" + strconv.Itoa(int(userID)))
  fmt.Println("redisValue for featureName:" + featureName)
  fmt.Println("feature:" + featureName + ",user:" + strconv.Itoa(int(userID)))
  fmt.Println(redisValue)
  
  return (redisValue == "x")
}
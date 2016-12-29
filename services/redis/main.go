package redis

import(
       "github.com/garyburd/redigo/redis"
       )

func GetByteArrayFromRedis(key string) []byte {
  client, err := redis.Dial("tcp", ":6379")
  if err != nil {
      panic(err)
  }
  defer client.Close()
  redisResponse, _ := redis.String(client.Do("GET", key))
  return []byte(redisResponse)
}


func SaveByteArrayToRedis(key string, value []byte) {
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


func SetRedisKey(key string, value string, ttl int) {
    client, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(err)
    }
    defer client.Close()
    client.Do("SET", key, value)
    client.Do("EXPIRE", key, ttl)
}

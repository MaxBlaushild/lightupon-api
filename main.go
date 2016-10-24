package main

import (
        "github.com/joho/godotenv"
        "lightupon-api/router"
        "lightupon-api/models"
        "lightupon-api/websockets"
        "lightupon-api/googleMaps"
        )

func main() {
  // load the environment variables from dotenv
  godotenv.Load()

  //connect to the database
  models.Connect()

  //connect to google maps api
  googleMaps.Init()

  //intialize the websocket hub and start waiting for connections
  go websockets.H.StartHub()

  //create the router and start listening for requests
  router.Init()
}

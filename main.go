package main

import (
        "github.com/joho/godotenv"
        "trip-advisor-backend/router"
        "trip-advisor-backend/models"
        "trip-advisor-backend/websockets"

        )

func main() {
  // load the environment variables from dotenv
  godotenv.Load()

  //connect to the database
  models.Connect()

  //intialize the websocket hub and start waiting for connections
  go websockets.H.StartHub()

  //create the router and start listening for requests
  router.Init()
}

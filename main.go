package main

import (
        "github.com/joho/godotenv"
        "lightupon-api/router"
        "lightupon-api/models"
        // "lightupon-api/live"
        "lightupon-api/services/twitter"
        )

func main() {
  // load the environment variables from dotenv
  godotenv.Load()

  //connect to the database
  models.Connect()

  // loads environment varaibles for twitter
  twitter.Init()

  //intialize the websocket hub and start waiting for connections
  // go live.Hub.Start()
  
  //create the router and start listening for requests
  router.Init()

}

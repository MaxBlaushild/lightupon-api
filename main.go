package main

import (
        "github.com/joho/godotenv"
        "lightupon-api/router"
        "lightupon-api/models"
        "lightupon-api/services/twitter"
        )

func main() {
  // load the environment variables from dotenv
  godotenv.Load()

  //connect to the database
  models.Connect(true)

  // loads environment varaibles for twitter
  twitter.Init()
  
  //create the router and start listening for requests
  router.Init()

}

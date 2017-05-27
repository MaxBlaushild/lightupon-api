package main

import (
        "github.com/joho/godotenv"
        "lightupon-api/router"
        "lightupon-api/models"
        "lightupon-api/live"
        )

func main() {
  // load the environment variables from dotenv
  godotenv.Load()

  //connect to the database
  models.Connect()

  //intialize the websocket hub and start waiting for connections
  go live.Hub.Start()
  // facebook.Post("EAANWMwWG4xABALd9yEAYWkfblFE0051PS2AspRSYjMPZBwYZCQWydxliYDbhxahGsSdvm4f80RE5SpCZCHFZCcFR7afutMJKSaZCa92IZAvRwULVZCurGMJ0U7355IoGFOVtzYwe7s5qqYUkmKzm768eoNyGC5tkNMEOx9yqTKacsjH3yRqZApYwBuibD61BXtRcXT1PrXZA2ncInjiQXKQT93X7NWuCpJicZD")
  //create the router and start listening for requests
  router.Init()

}

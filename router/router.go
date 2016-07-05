package router

import(
      "os"

      "github.com/gorilla/mux"
      "github.com/rs/cors"
      "github.com/codegangsta/negroni"

      "lightupon-api/router/routes"
      "lightupon-api/router/middleware"
      )

func Init(){
  muxRouter := mux.NewRouter().StrictSlash(true)
  muxRouter.HandleFunc("/lightupon/users", routes.UserLogisterHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/users/{facebookId}/token", routes.UserTokenRefreshHandler).Methods("PATCH")

  routerWithAuth := mux.NewRouter()
  
  routerWithAuth.HandleFunc("/lightupon/trips", routes.TripsHandler)
  routerWithAuth.HandleFunc("/lightupon/trips/{id}", routes.TripHandler)
  routerWithAuth.HandleFunc("/lightupon/nearby_trips", routes.NearbyTripsHandler)
  routerWithAuth.HandleFunc("/lightupon/trips/{tripId}/scenes", routes.ScenesHandler)
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneId}/cards", routes.CardsHandler)
  routerWithAuth.HandleFunc("/lightupon/parties", routes.CreatePartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties", routes.GetUsersPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{id}", routes.GetPartyHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/users", routes.AddUserToPartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/status", routes.UpdatePartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/pull", routes.PartyManagerHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/nextScene", routes.MovePartyToNextSceneHandler)
  routerWithAuth.HandleFunc("/lightupon/parties", routes.LeavePartyHandler).Methods("DELETE")
  // TODO: rename trips_post to trips and get it to not get confused with the above trips GET route
  routerWithAuth.HandleFunc("/lightupon/trips_post", routes.CreateTripHandler).Methods("POST")

  muxRouter.PathPrefix("/").Handler(negroni.New(
    negroni.HandlerFunc(middleware.Auth().HandlerWithNext),
    negroni.Wrap(routerWithAuth),
  ))

  port := os.Getenv("PORT")
  if (len(port) == 0) {
    port = "5000"
  }

  // apply CORS
  finalHandler := cors.Default().Handler(muxRouter)

  n := negroni.Classic()
  n.UseHandler(finalHandler)
  n.Run(":" + port)
}


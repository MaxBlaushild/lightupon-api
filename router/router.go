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

  // ADMIN
  // these routes serve html
  muxRouter.HandleFunc("/lightupon/admin/user/{id}/trips", routes.AdminGetTripsForUserHandler)
  muxRouter.HandleFunc("/lightupon/admin/trips/{id}", routes.AdminTripHandler)
  // these serve/accept json
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}/cards", routes.CardsHandler)
  muxRouter.HandleFunc("/lightupon/admin/trips/{tripID}/scenes_post", routes.CreateSceneHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}/cards_post", routes.CreateCardHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}", routes.DeleteSceneHandler).Methods("DELETE")
  muxRouter.HandleFunc("/lightupon/admin/cards/{cardID}", routes.DeleteCardHandler).Methods("DELETE")

  routerWithAuth := mux.NewRouter()
  
  routerWithAuth.HandleFunc("/lightupon/trips", routes.TripsHandler)
  routerWithAuth.HandleFunc("/lightupon/trips/{id}", routes.TripHandler)
  routerWithAuth.HandleFunc("/lightupon/trips_for_user", routes.GetTripsForUserHandler)
  routerWithAuth.HandleFunc("/lightupon/nearby_trips", routes.NearbyTripsHandler)
  routerWithAuth.HandleFunc("/lightupon/trips/{tripId}/scenes", routes.ScenesHandler)
  routerWithAuth.HandleFunc("/lightupon/trips/{tripID}/scenes_post", routes.CreateSceneHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/cards", routes.CardsHandler)
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/cards_post", routes.CreateCardHandler).Methods("POST")
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


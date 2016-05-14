package router

import(
      "os"

      "github.com/gorilla/mux"
      "github.com/rs/cors"
      "github.com/codegangsta/negroni"

      "trip-advisor-backend/router/routes"
      "trip-advisor-backend/router/middleware"
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
  routerWithAuth.HandleFunc("/lightupon/parties", routes.GetUsersPartyHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{id}", routes.GetPartyHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/users", routes.AddUserToPartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/status", routes.UpdatePartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyId}/users", routes.PartyMembersHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/pull", routes.PartyManagerHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/start", routes.StartPartyHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/leave", routes.LeavePartyHandler)

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
  n.Run(":5000")
}


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

  // HOMEPAGE
  muxRouter.HandleFunc("/lightupon/home/", routes.ServeHomepage).Methods("GET")

  // ADMIN
  // these routes serve html
  muxRouter.HandleFunc("/lightupon/admin/user/{id}/trips", routes.AdminGetTripsForUserHandler)
  muxRouter.HandleFunc("/lightupon/admin/trips/{id}", routes.AdminTripDetailsHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{id}", routes.AdminSceneDetailsHandler).Methods("GET")

  // these serve/accept json
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}/cards", routes.CardsHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/admin/popularScenes", routes.PopularScenesHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/admin/trips/{tripID}/scenes", routes.CreateSceneHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}/cards", routes.CreateCardHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/trips", routes.AdminCreateTripHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}", routes.DeleteSceneHandler).Methods("DELETE")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}", routes.ModifySceneHandler).Methods("PUT")
  muxRouter.HandleFunc("/lightupon/admin/cards/{cardID}", routes.ModifyCardHandler).Methods("PUT")
  muxRouter.HandleFunc("/lightupon/admin/cards/{cardID}", routes.DeleteCardHandler).Methods("DELETE")
  muxRouter.HandleFunc("/lightupon/admin/trips/{tripID}", routes.DeleteTripHandler).Methods("DELETE")
  muxRouter.HandleFunc("/lightupon/admin/assets/uploadUrls/{assetType}/{assetName}", routes.UploadAssetUrlHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/trips/", routes.CreateSelfieTripHandler).Methods("POST")

  routerWithAuth := mux.NewRouter()

  // USER STUFF
  routerWithAuth.HandleFunc("/lightupon/me", routes.MeHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/user/{userID}/follow", routes.FollowHandler).Methods("POST")

  // routerWithAuth.HandleFunc("/lightupon/getFollowers", routes.GetFollowersHandler).Methods("GET")
  // routerWithAuth.HandleFunc("/lightupon/getFolloweringUsers", routes.GetFollowingUsersHandler).Methods("GET")

  routerWithAuth.HandleFunc("/lightupon/users", routes.SearchUsersHandler).Methods("GET")
  
  // LIGHT STUFF
  routerWithAuth.HandleFunc("/lightupon/light", routes.LightHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/extinguish", routes.LightHandler).Methods("POST")
  
  // PARTY STUFF
  routerWithAuth.HandleFunc("/lightupon/trips", routes.CreateTripHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/trips", routes.TripsHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/tripsForUser", routes.GetTripsForUserHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/trips/{id}", routes.TripHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/trips/{tripId}/scenes", routes.ScenesHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/trips/{tripID}/scenes", routes.CreateSceneHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/cards", routes.CardsHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/cards", routes.CreateCardHandler).Methods("POST")
  
  // PARTY STUFF
  routerWithAuth.HandleFunc("/lightupon/parties", routes.GetUsersPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties", routes.CreatePartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{id}", routes.GetPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/users", routes.AddUserToPartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties", routes.LeavePartyHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/pull", routes.PullHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/finishParty", routes.FinishPartyHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/nextScene", routes.MovePartyToNextSceneHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/finishParty", routes.FinishPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/invite", routes.CreatePartyInviteHandler).Methods("POST")

  // BOOKMARKS
  muxRouter.HandleFunc("/lightupon/login/", routes.Login).Methods("GET")
  muxRouter.HandleFunc("/lightupon/bookmarks/", routes.ServeBookmarks).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/me/bookmarks", routes.GetBookmarksForUser).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/bookmarks/{bookmarkID}/like", routes.LikeBookmark).Methods("PUT")
  routerWithAuth.HandleFunc("/lightupon/bookmarks/{bookmarkID}/fuckThis", routes.FuckThisBookmark).Methods("PUT")

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


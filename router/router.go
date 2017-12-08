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

  routerWithAuth := mux.NewRouter()

  // USER STUFF
  routerWithAuth.HandleFunc("/lightupon/me", routes.MeHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/follow", routes.FollowHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/follow", routes.UnfollowHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/users", routes.SearchUsersHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}", routes.GetUserHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/deviceToken", routes.AddDeviceToken).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/me/twitter/login", routes.TwitterLoginHandler).Methods("POST")


  // // LOCATION STUFF
  routerWithAuth.HandleFunc("/lightupon/discover", routes.DiscoverHandler).Methods("POST")
  
  // PARTY STUFF
  routerWithAuth.HandleFunc("/lightupon/parties", routes.GetUsersPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties", routes.CreatePartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties/{id}", routes.GetPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{passcode}/users", routes.AddUserToPartyHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/parties", routes.LeavePartyHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/pull", routes.PullHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/nextScene", routes.MovePartyToNextSceneHandler)
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/end", routes.EndPartyHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/parties/{partyID}/invite", routes.CreatePartyInviteHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/admin/assets/uploadUrls", routes.UploadAssetUrlHandler).Methods("POST")

  // POSTS STUFF
  routerWithAuth.HandleFunc("/lightupon/posts", routes.CreatePost).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/posts/{postID}", routes.GetPostHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/posts", routes.GetPostHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/posts/nearby", routes.GetNearbyPosts).Methods("GET")

  // VOTES
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/upvote", routes.PostUpvoteHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/downvote", routes.PostDownvoteHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/vote", routes.DeleteVoteHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/voteTotal", routes.GetVoteTotalHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/user/walletTotal", routes.GetWalletPointsHandler).Methods("GET")

  // WEB STUFF
  muxRouter.HandleFunc("/lightupon/login/", routes.Login).Methods("GET")
  muxRouter.HandleFunc("/lightupon/stats/", routes.ServeStatsPage).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/stats/json", routes.GetStats).Methods("GET")

  muxRouter.PathPrefix("/").Handler(negroni.New(
    negroni.HandlerFunc(middleware.Auth().HandlerWithNext),
    negroni.Wrap(routerWithAuth),
  ))

  port := os.Getenv("PORT")
  if (len(port) == 0) {
    port = "5000"
  }

  // apply CORS

  c := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
    AllowCredentials: true,
  })

  finalHandler := c.Handler(muxRouter)

  n := negroni.Classic()
  n.UseHandler(finalHandler)
  n.Run(":" + port)
}


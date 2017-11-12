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
  muxRouter.HandleFunc("/lightupon/admin/trips/{tripID}/scenes", routes.CreateSceneHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}", routes.DeleteSceneHandler).Methods("DELETE")
  muxRouter.HandleFunc("/lightupon/admin/scenes/{sceneID}", routes.ModifySceneHandler).Methods("PUT")
  // muxRouter.HandleFunc("/lightupon/me/instagram/login", routes.InstagramLoginHandler).Methods("GET")

  routerWithAuth := mux.NewRouter()

  // USER STUFF
  routerWithAuth.HandleFunc("/lightupon/me", routes.MeHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/follow", routes.FollowHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/follow", routes.UnfollowHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/feed", routes.FollowingScenesHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users", routes.SearchUsersHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}", routes.GetUserHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/deviceToken", routes.AddDeviceToken).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/me/twitter/login", routes.TwitterLoginHandler).Methods("POST")

  // LIGHT STUFF
  routerWithAuth.HandleFunc("/lightupon/light", routes.LightHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/extinguish", routes.ExtinguishHandler).Methods("POST")
  
  // TRIP STUFF
  routerWithAuth.HandleFunc("/lightupon/scenes/nearby", routes.NearbyScenesHandler).Methods("GET")
  // routerWithAuth.HandleFunc("/lightupon/scenes/nearby", routes.NearbyFollowingScenesHandler).Methods("GET")

  // LOCATION STUFF
  routerWithAuth.HandleFunc("/lightupon/discover", routes.DiscoverHandler).Methods("POST")

  // SCENE STUFF
  routerWithAuth.HandleFunc("/lightupon/scenes", routes.ScenesIndexHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/activeScene", routes.ActiveSceneHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/scenes", routes.ScenesForUserHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/discover", routes.DiscoverSceneHandler).Methods("POST")
  
  // PARTY STUFF
  routerWithAuth.HandleFunc("/lightupon/pull", routes.PullHandler).Methods("GET")

  // COMMENTS STUFF
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}", routes.SceneHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/comments", routes.SceneCommentsHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/comments", routes.PostSceneCommentHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/likes", routes.LikeSceneHandler).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/likes", routes.UnlikeSceneHandler).Methods("DELETE")
  routerWithAuth.HandleFunc("/lightupon/scenes/{sceneID}/flag", routes.FlagSceneHandler).Methods("POST")


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
    port = "5001"
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


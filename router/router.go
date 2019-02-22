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
  muxRouter.HandleFunc("/health", routes.HealthHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/users", routes.UserLogisterHandler).Methods("POST")
  muxRouter.HandleFunc("/lightupon/users/{facebookId}/token", routes.UserTokenRefreshHandler).Methods("PATCH")

  routerWithAuth := mux.NewRouter()

  // USER STUFF
  muxRouter.HandleFunc("/lightupon/me", routes.MeHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/users", routes.SearchUsersHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/users/{userID}", routes.GetUserHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/deviceToken", routes.AddDeviceToken).Methods("POST")
  muxRouter.HandleFunc("/lightupon/me/twitter/login", routes.TwitterLoginHandler).Methods("POST")

  // PARTY STUFF
  muxRouter.HandleFunc("/lightupon/admin/assets/uploadUrls", routes.UploadAssetUrlHandler).Methods("POST")

  // POSTS STUFF
  muxRouter.HandleFunc("/lightupon/posts", routes.CreatePost).Methods("POST")
  muxRouter.HandleFunc("/lightupon/users/{userID}/posts", routes.GetUsersPosts).Methods("GET")
  muxRouter.HandleFunc("/lightupon/posts/nearby", routes.GetNearbyPostsAndTryToDiscoverThem).Methods("GET")
  muxRouter.HandleFunc("/lightupon/posts/{postID}", routes.GetPostHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/posts/{postID}/complete", routes.CompletePostHandler).Methods("POST")

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


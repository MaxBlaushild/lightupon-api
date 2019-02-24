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

  // USERS
  routerWithAuth.HandleFunc("/lightupon/me", routes.MeHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users", routes.SearchUsersHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}", routes.GetUserHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/deviceToken", routes.AddDeviceToken).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/me/twitter/login", routes.TwitterLoginHandler).Methods("POST")

  // POSTS
  routerWithAuth.HandleFunc("/lightupon/posts", routes.CreatePost).Methods("POST")
  routerWithAuth.HandleFunc("/lightupon/users/{userID}/posts", routes.GetUsersPosts).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/posts/nearby", routes.GetNearbyPostsAndTryToDiscoverThem).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/posts/{postID}", routes.GetPostHandler).Methods("GET")
  routerWithAuth.HandleFunc("/lightupon/posts/{postID}/complete", routes.CompletePostHandler).Methods("POST")

  // QUESTS
  muxRouter.HandleFunc("/lightupon/quests", routes.AllQuestsHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/quests/{questID}/edit", routes.EditQuestHandler).Methods("GET")
  muxRouter.HandleFunc("/lightupon/quests/{questID}/update", routes.UpdateQuestHandler).Methods("POST")

  routerWithAuth.HandleFunc("/lightupon/admin/assets/uploadUrls", routes.UploadAssetUrlHandler).Methods("POST")

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
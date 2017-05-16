package routes

import(
       "net/http"
       )

func ServeStatsPage(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/stats.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "html/login.html")
}
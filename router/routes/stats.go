package routes

import(
       "net/http"
       "lightupon-api/models"
       "encoding/json"
       )


func GetStats(w http.ResponseWriter, r *http.Request) {
  stats := models.GetUserStats()
  json.NewEncoder(w).Encode(stats)
}


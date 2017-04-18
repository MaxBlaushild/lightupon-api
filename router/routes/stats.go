package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/davecgh/go-spew/spew"
       "fmt"
       "encoding/json"
       // "strconv"

       )


func GetStats(w http.ResponseWriter, r *http.Request) {

  stats := models.GetUserStats()
  fmt.Println("spew.Dump(stats)")
  spew.Dump(stats)

  json.NewEncoder(w).Encode(stats)
}


package routes

import(
       "net/http"
       "lightupon-api/live"
       "fmt"
)

func PullHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  // activeParty := user.ActiveParty()

  ws, err := live.Upgrader.Upgrade(w, r, nil); if err != nil {
    return
  }

  c := &live.Connection{
    Send: make(chan live.Response), 
    WS: ws, 
    UserID: user.ID,
  }

  // activeParty.Connect(c)
  fmt.Println(c)
}
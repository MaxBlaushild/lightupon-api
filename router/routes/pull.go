package routes

import(
       "net/http"
       "lightupon-api/models"
       "lightupon-api/websockets"
       "fmt"
)

func PullHandler(w http.ResponseWriter, r *http.Request) {
  user := GetUserFromRequest(r)
  activeParty := user.ActiveParty()
  ws, err := websockets.Upgrader.Upgrade(w, r, nil); if err != nil {
    respondWithBadRequest(w, "You done fucked up. Give us a real passcode.")
    return
  }

  fmt.Println(activeParty)

  c := &websockets.Connection{Send: make(chan models.PullResponse), Passcode: activeParty.Passcode, WS: ws, User: user}

  websockets.H.Register <- c

  go c.ReadPump()
  c.WritePump()
}
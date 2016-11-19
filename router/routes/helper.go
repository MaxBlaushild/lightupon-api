package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/context"
       "github.com/dgrijalva/jwt-go"
       "github.com/gorilla/mux"
       "strconv"
       )

func GetUserFromRequest(r *http.Request)(user models.User){
  token := context.Get(r, "user")
  facebookID := token.(*jwt.Token).Claims["facebookId"].(string)
  models.DB.Preload("Follows").Where("facebook_id = ?", facebookID).First(&user)
  return
}

func GetUIntFromVars(r *http.Request, field string)(uintToReturn uint){
  vars := mux.Vars(r)
  george, _ := strconv.Atoi(vars[field])
  uintToReturn = uint(george) // Fuck unints
  return
}

func GetUserLocationFromRequest(r *http.Request)(lat string, lon string){
  query := r.URL.Query()
  lat = query["lat"][0]
  lon = query["lon"][0]
  return
}
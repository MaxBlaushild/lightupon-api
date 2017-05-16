package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/gorilla/context"
       "github.com/dgrijalva/jwt-go"
       "github.com/gorilla/mux"
       "strconv"
       )

func GetFacebookIdFromRequest(r *http.Request) string {
  token := context.Get(r, "user")
  facebookID := token.(*jwt.Token).Claims["facebookId"].(string)
  return facebookID
}

func GetUserFromRequest(r *http.Request) (user models.User) {
  facebookID := GetFacebookIdFromRequest(r)
  models.DB.Where("facebook_id = ?", facebookID).First(&user)
  return
}

func GetUIntFromVars(r *http.Request, field string) (uint, error) {
  vars := mux.Vars(r)
  intValue, err := strconv.Atoi(vars[field])
  uintValue := uint(intValue) // Fuck unints
  return uintValue, err
}

func GetUserLocationFromRequest(r *http.Request) (lat string, lon string) {
  query := r.URL.Query()
  lat = query["lat"][0]
  lon = query["lon"][0]
  return
}

func getStringFromRequest(r *http.Request, key string, defaultValue string) (value string) {
  query := r.URL.Query()
  if len(query[key]) > 0 {
    value = query[key][0]
  } else {
    value = defaultValue
  }
  return
}
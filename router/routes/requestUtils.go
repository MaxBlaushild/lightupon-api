package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/dgrijalva/jwt-go"
       "github.com/gorilla/mux"
       "strconv"
       )

func GetLocationFromRequest(r *http.Request) (lat string, lon string) {
  query := r.URL.Query()
  lat = query["lat"][0]
  lon = query["lon"][0]
  return
}

func GetUserFromRequest(r *http.Request) (user models.User) {
  facebookID := getFacebookIdFromRequest(r)
  models.DB.Where("facebook_id = ?", facebookID).First(&user)
  return
}

func GetStringFromRequest(r *http.Request, key string, defaultValue string) (value string) {
  query := r.URL.Query()
  if len(query[key]) > 0 {
    value = query[key][0]
  } else {
    value = defaultValue
  }
  return
}

func GetUIntFromVars(r *http.Request, field string) (uint, error) {
  vars := mux.Vars(r)
  intValue, err := strconv.Atoi(vars[field])
  uintValue := uint(intValue)
  return uintValue, err
}

func getFacebookIdFromRequest(r *http.Request) string {
  ctx := r.Context()
  token := ctx.Value("user").(*jwt.Token)
  claims := token.Claims.(jwt.MapClaims)
  facebookID := claims["facebookId"].(string)
  return facebookID
}
package routes

import(
       "net/http"
       "lightupon-api/models"
       "github.com/dgrijalva/jwt-go"
       "github.com/gorilla/mux"
       "strconv"
       )

/* 
Taking an object-oriented approach here so we can use dependency injection.
*/

type requestAccessor interface {
  GetLocationFromRequest() (lat string, lon string)
  GetUserFromRequest() (user models.User)
  getStringFromRequest(key string, defaultValue string) (value string)
}

type requestManager struct {
  request *http.Request
}

func newRequestManager(request *http.Request) (requestManager requestManager) {
  requestManager.request = request
  return
}

func (rm requestManager) GetLocationFromRequest() (lat string, lon string) {
  query := rm.request.URL.Query()
  lat = query["lat"][0]
  lon = query["lon"][0]
  return
}

func (rm requestManager) GetUserFromRequest() (user models.User) {
  facebookID := getFacebookIdFromRequest(rm.request)
  models.DB.Where("facebook_id = ?", facebookID).First(&user)
  return
}

func (rm requestManager) getStringFromRequest(key string, defaultValue string) (value string) {
  query := rm.request.URL.Query()
  if len(query[key]) > 0 {
    value = query[key][0]
  } else {
    value = defaultValue
  }
  return
}

func (rm requestManager) GetUIntFromVars(field string) (uint, error) {
  vars := mux.Vars(rm.request)
  intValue, err := strconv.Atoi(vars[field])
  uintValue := uint(intValue) // Fuck unints
  return uintValue, err
}

func getFacebookIdFromRequest(r *http.Request) string {
  ctx := r.Context()
  token := ctx.Value("user").(*jwt.Token)
  claims := token.Claims.(jwt.MapClaims)
  facebookID := claims["facebookId"].(string)
  return facebookID
}
package middleware

import(
			"os"
      "fmt"
      "net/http"
			
    	"github.com/auth0/go-jwt-middleware"
    	"github.com/dgrijalva/jwt-go"
)

func OnError(w http.ResponseWriter, r *http.Request, err string) {
  http.Error(w, err, http.StatusUnauthorized)
}

func Auth() *jwtmiddleware.JWTMiddleware {
  jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      return []byte(os.Getenv("JWT_SECRET")), nil
    },
    ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
      fmt.Println(err)
      http.Error(w, err, http.StatusUnauthorized)
    },
    SigningMethod: jwt.SigningMethodHS256,
  })
  return jwtMiddleware
}
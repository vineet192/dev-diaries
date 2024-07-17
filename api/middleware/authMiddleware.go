package middleware

import (
	"devdiaries/api/utilities"
	secretsvault "devdiaries/secrets_vault"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtSecret string
var token *jwt.Token

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if len(authHeader) == 0 {
			http.Error(w, "No Authorization header present", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		//refresh should occur only once per runtime
		if len(jwtSecret) == 0 {
			if err := refreshSecret(); err != nil {
				http.Error(w,
					"Unexpected error retrieving server credentials",
					http.StatusInternalServerError)
				return
			}
		}

		var jwtParseErr error

		token, jwtParseErr = jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if jwtParseErr != nil {
			utilities.HandleJWTError(jwtParseErr, r.URL.String(), w)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func ValidateUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		tokenID := token.Claims.(*jwt.RegisteredClaims).ID

		if id != tokenID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func refreshSecret() (err error) {
	jwtSecret, err = secretsvault.GetSecret("JWT_SECRET")

	if err != nil {
		return err
	}
	return
}

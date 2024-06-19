package auth

import (
	"devdiaries/api/utilities"
	"devdiaries/database"
	"devdiaries/models"
	"devdiaries/payload/request"
	secretsvault "devdiaries/secrets_vault"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret string

func Login(w http.ResponseWriter, r *http.Request) {
	var loginBody request.LoginBody
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&loginBody); err != nil {
		utilities.HandleJSONDecodeErr(err, r.URL.String(), w)
	}

	if err := database.DB.Where(&models.User{Email: loginBody.Email}).First(&user).Error; err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(loginBody.Password))

	if hashErr != nil {
		utilities.HandleHashError(hashErr, r.URL.String(), w)
		return
	}

	ss, jwtErr := mintJWTToken(user.ID)

	if jwtErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createCookie(w, ss)
	w.WriteHeader(http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var signupBody request.SignupBody

	if err := json.NewDecoder(r.Body).Decode(&signupBody); err != nil {
		utilities.HandleJSONDecodeErr(err, r.URL.String(), w)
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(signupBody.Password), 14)

	if hashErr != nil {
		utilities.HandleHashError(hashErr, r.URL.String(), w)
		return
	}

	signupBody.User.Hash = string(hash)

	tx := database.DB.Begin()
	dbErr := createUser(&signupBody.User, tx)

	if dbErr != nil {
		tx.Rollback()
		utilities.HandleDBError(dbErr, r.URL.String(), w, "user")
		return
	}

	ss, jwtErr := mintJWTToken(signupBody.User.ID)

	if jwtErr != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createCookie(w, ss)

	tx.Commit()
	w.WriteHeader(http.StatusCreated)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func createUser(user *models.User, tx *gorm.DB) error {

	result := tx.Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func createCookie(w http.ResponseWriter, ss string) {

	minutesOffset, err := strconv.ParseUint(os.Getenv("TOKEN_VALIDITY_MINUTES"), 10, 64)

	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Expires: time.Now().Add(time.Minute * time.Duration(minutesOffset)),
		Value:   ss,
		Name:    "devdiaries_user",
	}
	http.SetCookie(w, &cookie)
}

func mintJWTToken(id uint) (ss string, err error) {

	if jwtSecret == "" {
		jwtSecret, err = secretsvault.GetSecret("JWT_SECRET")
	}

	if err != nil {
		return "", err
	}

	minutesOffset, err := strconv.ParseUint(os.Getenv("TOKEN_VALIDITY_MINUTES"), 10, 64)

	if err != nil {
		panic(err)
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(minutesOffset))),
		Issuer:    "devdiaries",
		ID:        fmt.Sprint(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return ss, nil
}

package utilities

import (
	"devdiaries/payload/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var foreignKeyErrorMap map[string]string = map[string]string{
	"blog":     "author does not exist or invalid tag",
	"comment":  "user or blog does not exist",
	"reaction": "blog or comment does not exist",
	"tag":      "tag or blog does not exist",
}

var duplicateKeyErrorMap map[string]string = map[string]string{
	"blog":     "Duplicated key",
	"comment":  "Duplicated key",
	"reaction": "Duplicated key",
	"tag":      "Duplicated key",
	"user":     "This user/email has already been registered",
}

func HandleJSONDecodeErr(err error, instance string, w http.ResponseWriter) {

	var msg string
	var detail string

	var syntaxErr *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var unsupportedTypeError *json.UnsupportedTypeError

	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.As(err, &syntaxErr):
		msg = "JSON syntax error"
		detail = err.Error()
		w.WriteHeader(http.StatusBadRequest)

	case errors.As(err, &unmarshalTypeError):
		msg = "Unexpected fields"
		detail = fmt.Sprintf(`The %q field %q was not expected`, unmarshalTypeError.Type, unmarshalTypeError.Field)
		w.WriteHeader(http.StatusBadRequest)

	case errors.As(err, &unsupportedTypeError):
		msg = "Unsupported type in JSON"
		detail = err.Error()
		w.WriteHeader(http.StatusBadRequest)

	case errors.Is(err, io.EOF):
		msg = "Empty body"
		detail = "Empty body"
		w.WriteHeader(http.StatusBadRequest)

	case err.Error() == "http: request body too large":
		msg = "Request body too large"
		detail = "Request body must not exceed 1MB"
		w.WriteHeader(http.StatusRequestEntityTooLarge)

	case err.Error() == "unexpected EOF":
		msg = "unexpected EOF"
		detail = "unexptected EOF"
		w.WriteHeader(http.StatusBadRequest)

	default:
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	errorResp, _ := json.Marshal(response.ErrorResponse{
		Message:  msg,
		Detail:   detail,
		Instance: instance,
	})

	w.Write(errorResp)
}

func HandleDBError(err error, instance string, w http.ResponseWriter, entity string) {

	var detail string
	var msg string

	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		w.WriteHeader(http.StatusBadRequest)
		msg = "missing constraint"
		detail = foreignKeyErrorMap[entity]
	case errors.Is(err, gorm.ErrInvalidValue):
		w.WriteHeader(http.StatusBadRequest)
		msg = "invalid value"
		detail = "invalid value"
	case errors.Is(err, gorm.ErrRecordNotFound):
		w.WriteHeader(http.StatusNotFound)
		msg = "record not found"
		detail = "record not found"
	case errors.Is(err, gorm.ErrDuplicatedKey):
		w.WriteHeader(http.StatusForbidden)
		msg = "Duplicated key"
		detail = duplicateKeyErrorMap[entity]
	default:
		fmt.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	errorResp, _ := json.Marshal(response.ErrorResponse{
		Message:  msg,
		Detail:   detail,
		Instance: instance,
	})

	w.Write(errorResp)
}

func HandleHashError(err error, instance string, w http.ResponseWriter) {
	var detail string
	var msg string

	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, bcrypt.ErrPasswordTooLong):
		w.WriteHeader(http.StatusBadRequest)
		msg = "Password too long"
		detail = "Password cannot exceed 72 bytes"
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		w.WriteHeader(http.StatusUnauthorized)
		msg = "Invalid Password"
		detail = "Invalid Password"
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errorResp, _ := json.Marshal(response.ErrorResponse{
		Message:  msg,
		Detail:   detail,
		Instance: instance,
	})

	w.Write(errorResp)
}

func HandleJWTError(err error, instance string, w http.ResponseWriter) {
	var detail string
	var msg string

	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		msg = "Malformed auth token"
		detail = "That's not even a token"
		w.WriteHeader(http.StatusBadRequest)

	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		msg = "Invalid signature"
		detail = "Invalid signature"
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		msg = "Token expired or not valid yet"
		detail = "Token expired or not valid yet"
		w.WriteHeader(http.StatusUnauthorized)
	default:
		msg = "Unexpected error"
		detail = "Unexpected error while parsing token"
		w.WriteHeader(http.StatusInternalServerError)
	}

	errorResp, _ := json.Marshal(response.ErrorResponse{
		Message:  msg,
		Detail:   detail,
		Instance: instance,
	})
	w.Write(errorResp)

}

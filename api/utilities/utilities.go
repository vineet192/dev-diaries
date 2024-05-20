package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"inventory/payload/response"
	"io"
	"net/http"

	"gorm.io/gorm"
)

func HandleJSONDecodeErr(err error, instance string, w http.ResponseWriter) {

	var msg string
	var detail string

	var syntaxErr *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var unsupportedTypeError *json.UnsupportedTypeError

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

func HandleDBErrorOnCreateBlog(err error, instance string, w http.ResponseWriter) {

	var detail string
	var msg string

	switch {
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		w.WriteHeader(http.StatusBadRequest)
		msg = "missing constraint"
		detail = "author does not exist"
	case errors.Is(err, gorm.ErrInvalidValue):
		w.WriteHeader(http.StatusBadRequest)
		msg = "invalid value"
		detail = "invalid value"
	default:
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

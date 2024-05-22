package user

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"io"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	decodeErr := decodeJSONBody(&user, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "user")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decodeErr := decodeJSONBody(&user, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.UpdateColumns(user)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "user")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	var user models.User

	decodeErr := decodeJSONBody(&user, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
	}

	result := database.DB.Delete(user)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "user")
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeJSONBody(blog *models.User, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(blog)
}

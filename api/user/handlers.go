package user

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"inventory/payload/response"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decodeErr := dec.Decode(&user)
	w.Header().Add("Content-Type", "application/json")

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		errorResp := response.ErrorResponse{
			Message:  "Error inserting into database",
			Instance: r.URL.String(),
			Detail:   result.Error.Error()}

		jsonErrorResp, jsonErr := json.Marshal(errorResp)

		if jsonErr != nil {
			panic(jsonErr)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErrorResp)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

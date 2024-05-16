package api

import (
	"encoding/json"
	"inventory/database"
	"inventory/models"
	"inventory/payload/response"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	jsonDecodeErr := json.NewDecoder(r.Body).Decode(&user)

	if jsonDecodeErr != nil {
		errorResp := createErrorResp("Error decoding JSON", r.URL.String(), jsonDecodeErr.Error())
		jsonErrorResp, jsonErr := json.Marshal(errorResp)

		if jsonErr != nil {
			panic(jsonErr)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErrorResp)
		return
	}

	result := database.DB.Create(&user) // pass pointer of data to Create

	if result.Error != nil {
		errorResp := createErrorResp("Error inserting into database", r.URL.String(), result.Error.Error())
		jsonErrorResp, jsonErr := json.Marshal(errorResp)

		if jsonErr != nil {
			panic(jsonErr)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErrorResp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		RowsAffected int64  `json:"rows_affected"`
		Msg          string `json:"msg"`
	}{
		RowsAffected: result.RowsAffected,
		Msg:          "Created Successfully",
	})
}

func createErrorResp(msg string, instance string, detail string) response.ErrorResponse {
	return response.ErrorResponse{
		Message:  msg,
		Instance: instance,
		Detail:   detail,
	}
}

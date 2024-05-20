package blog

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"net/http"
)

func PostBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decodeErr := dec.Decode(&blog)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.Create(&blog)

	if result.Error != nil {

		utilities.HandleDBErrorOnCreateBlog(result.Error, r.URL.String(), w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

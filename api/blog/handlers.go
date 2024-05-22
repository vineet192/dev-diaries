package blog

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"inventory/payload/response"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PostBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog

	decodeErr := decodeJSONBody(&blog, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.Create(&blog)

	if result.Error != nil {

		utilities.HandleDBError(result.Error, r.URL.String(), w, "blog")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func EditBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog

	decodeErr := decodeJSONBody(&blog, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.UpdateColumns(blog)

	if result.Error != nil {

		utilities.HandleDBError(result.Error, r.URL.String(), w, "blog")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBlogByID(w http.ResponseWriter, r *http.Request) {

	id, parseErr := strconv.Atoi(mux.Vars(r)["id"])

	if parseErr != nil {
		errorResp, err := json.Marshal(response.ErrorResponse{
			Message:  "Invalid ID",
			Detail:   "Invalid ID",
			Instance: r.URL.String()})

		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(errorResp)
		return
	}

	result := database.DB.Delete(&models.Blog{}, id)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "blog")
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeJSONBody(blog *models.Blog, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(blog)
}

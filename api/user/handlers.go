package user

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"inventory/payload/request"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
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

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	var blogs []models.Blog
	var user models.User
	var params request.BlogQuery

	queries := r.URL.Query()

	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.ID = uint(id)

	decodeBlogQuery(&queries, &params)

	dbErr := executeBlogQuery(&user, &params, &blogs)

	if dbErr != nil {
		utilities.HandleDBError(dbErr, r.URL.String(), w, "user")
		return
	}

	response, _ := json.Marshal(blogs)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.Find(&user, id)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "user")
		return
	}

	resp, jsonMarshallErr := json.Marshal(user)

	if jsonMarshallErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(resp)

}

func decodeJSONBody(blog *models.User, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(blog)
}

func decodeBlogQuery(queries *url.Values, params *request.BlogQuery) error {
	decoder := schema.NewDecoder()

	return decoder.Decode(params, *queries)
}

func executeBlogQuery(user *models.User, params *request.BlogQuery, blogs *[]models.Blog) (err error) {
	var dbErr error

	query, values := params.DBQuery()

	db := database.DB.Model(user)
	db = db.Preload("Tags")

	if !params.DisableComments {
		db = db.Preload("Comments")
	}

	if !params.DisableReactions {
		db = db.Preload("Reactions")
	}

	if len(values) > 0 {
		dbErr = db.Where(query, values...).Scopes(Paginate(params)).Association("Blogs").Find(&blogs)
	} else {
		dbErr = db.Scopes(Paginate(params)).Association("Blogs").Find(&blogs)
	}

	if dbErr != nil {
		return dbErr
	}
	return
}

func Paginate(bq *request.BlogQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := bq.Page
		if page <= 0 {
			page = 1
		}

		pageSize := bq.PageSize

		switch {
		case pageSize > 20:
			pageSize = 20
		case pageSize == 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

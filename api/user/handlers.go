package user

import (
	"devdiaries/api/utilities"
	"devdiaries/database"
	"devdiaries/models"
	"devdiaries/payload/request"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decodeErr := decodeJSONBody(&user, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	user.ID = uint(id)

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

func GetBlogFeed(w http.ResponseWriter, r *http.Request) {
	var params request.BlogQuery
	var following []models.User
	var blogs []models.Blog
	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	blogQuery := r.URL.Query()

	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decodeBlogQuery(&blogQuery, &params)

	if err := database.DB.
		Model(&models.User{ID: uint(id)}).
		Association("Following").
		Find(&following); err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	if err := database.DB.
		Model(&models.Blog{}).
		Scopes(Paginate(&params)).
		Where("author_id IN ?", UserIDs(&following)).
		Order("posted_on desc").
		Find(&blogs).Error; err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	response, _ := json.Marshal(blogs)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func GetUserByID(w http.ResponseWriter, r *http.Request) {

	var user models.User

	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.First(&user, id)

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

func GetUser(w http.ResponseWriter, r *http.Request) {
	var params request.UserQuery
	var user models.User

	query := r.URL.Query()

	if err := decodeUserQuery(&query, &params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := executeUserQuery(&params, &user); err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
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

func AddFollower(w http.ResponseWriter, r *http.Request) {
	id, parseIDErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	follower_id, parseFollowerIDErr := strconv.ParseUint(mux.Vars(r)["follower_id"], 10, 64)

	var user models.User
	var follower models.User

	if parseIDErr != nil || parseFollowerIDErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := database.DB.Begin()

	if err := tx.First(&user, id).Error; err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	if err := tx.First(&follower, follower_id).Error; err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	if err := tx.
		Model(&user).
		Association("Followers").
		Append(&follower); err != nil {
		tx.Rollback()
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utilities.HandleDBError(err, r.URL.String(), w, "user")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func RemoveFollower(w http.ResponseWriter, r *http.Request) {
	id, parseIDErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	follower_id, parseFollowerIDErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseIDErr != nil || parseFollowerIDErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{ID: uint(id)}
	follower := models.User{ID: uint(follower_id)}

	result := database.DB.Model(&user).Delete(follower)

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

func decodeBlogQuery(queries *url.Values, params *request.BlogQuery) error {
	decoder := schema.NewDecoder()

	return decoder.Decode(params, *queries)
}

func decodeUserQuery(queries *url.Values, params *request.UserQuery) error {
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

func executeUserQuery(params *request.UserQuery, user *models.User) (err error) {

	db := database.DB.Model(&models.User{})

	if params.IncludeBlogReactions {
		db = db.Preload("BlogReactions")
	}

	if params.IncludeBlogs {
		db = db.Preload("Blogs")
	}

	if params.IncludeCommentReactions {
		db = db.Preload("CommentReactions")
	}

	if params.IncludeComments {
		db = db.Preload("Comments")
	}

	if params.IncludeFollowers {
		db = db.Preload("Followers")
	}

	if params.IncludeFollowing {
		db = db.Preload("Following")
	}

	if err := db.Where(&params.User).First(&user).Error; err != nil {
		return err
	}

	return nil

}

func UserIDs(users *[]models.User) []int {
	userIDs := make([]int, len(*users))

	for _, user := range *users {
		userIDs = append(userIDs, int(user.ID))
	}

	return userIDs
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

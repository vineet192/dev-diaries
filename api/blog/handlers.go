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

	decodeErr := decodeBlogJSONBody(&blog, &r.Body, &w)

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

func PostComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment

	blogID, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	decodeErr := decodeCommentJSONBody(&comment, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	comment.BlogID = uint(blogID)

	result := database.DB.Create(&comment)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "comment")
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func PostReaction(w http.ResponseWriter, r *http.Request) {
	var reaction models.BlogReaction

	blogID, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	decodeErr := decodeReactionJSONBody(&reaction, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	reaction.BlogID = uint(blogID)

	result := database.DB.Create(&reaction)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "reaction")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteReaction(w http.ResponseWriter, r *http.Request) {
	blog_id, blogIDParseErr := strconv.ParseUint(mux.Vars(r)["blog_id"], 10, 64)
	user_id, userIDParseErr := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	if blogIDParseErr != nil || userIDParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&models.BlogReaction{BlogID: uint(blog_id), UserID: uint(user_id)})

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "blog_reaction")
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func EditBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog

	id, parseErr := strconv.Atoi(mux.Vars(r)["id"])

	if parseErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	decodeErr := decodeBlogJSONBody(&blog, &r.Body, &w)
	blog.ID = uint(id)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	result := database.DB.Omit("Tags.*", "Reactions.*", "Comments.*").Model(&blog).Updates(blog)

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

func decodeBlogJSONBody(blog *models.Blog, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(blog)
}

func decodeCommentJSONBody(comment *models.Comment, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(comment)
}

func decodeReactionJSONBody(reaction *models.BlogReaction, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(reaction)
}

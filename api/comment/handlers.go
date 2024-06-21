package comment

import (
	"devdiaries/api/utilities"
	"devdiaries/database"
	"devdiaries/models"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Deletes a comment identifies by id
func DeleteCommentByID(w http.ResponseWriter, r *http.Request) {

	id, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := database.DB.Delete(&models.Comment{}, id)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "comment")
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// PostReaction accepts a CommentReaction in the request body and
// creates it against the comment identified by the id url parameter
func PostReaction(w http.ResponseWriter, r *http.Request) {
	var reaction models.CommentReaction

	commentID, parseErr := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)

	if parseErr != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	decodeErr := decodeReactionJSONBody(&reaction, &r.Body, &w)

	if decodeErr != nil {
		utilities.HandleJSONDecodeErr(decodeErr, r.URL.String(), w)
		return
	}

	reaction.CommentID = uint(commentID)

	result := database.DB.Create(&reaction)

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "comment_reaction")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Deletes a reaction posted by a user identified by the user_id url parameter for
// the comment identified by the comment_id url parameter
func DeleteReaction(w http.ResponseWriter, r *http.Request) {
	comment_id, blogIDParseErr := strconv.ParseUint(mux.Vars(r)["comment_id"], 10, 64)
	user_id, userIDParseErr := strconv.ParseUint(mux.Vars(r)["user_id"], 10, 64)

	if blogIDParseErr != nil || userIDParseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&models.CommentReaction{CommentID: uint(comment_id), UserID: uint(user_id)})

	if result.Error != nil {
		utilities.HandleDBError(result.Error, r.URL.String(), w, "comment_reaction")
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Decodes a request body into a CommentReaction object
func decodeReactionJSONBody(reaction *models.CommentReaction, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(reaction)
}

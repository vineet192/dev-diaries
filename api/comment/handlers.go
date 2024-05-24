package comment

import (
	"encoding/json"
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func decodeReactionJSONBody(reaction *models.CommentReaction, body *io.ReadCloser, w *http.ResponseWriter) error {
	dec := json.NewDecoder(*body)
	dec.DisallowUnknownFields()
	*body = http.MaxBytesReader(*w, *body, 1048576)

	return dec.Decode(reaction)
}

package comment

import (
	"inventory/api/utilities"
	"inventory/database"
	"inventory/models"
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

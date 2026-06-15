package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) deleteEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "DELETE", false) {
		return
	}

	id := r.PathValue("id")
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	err := cfg.deleteUserEntry(id, user_id)
	if err != nil {
		if err.Error() == "not found or unauthorized" {
			http.Error(w, "Transaction not found or unauthorized", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Transaction deleted successfully"}`))
}

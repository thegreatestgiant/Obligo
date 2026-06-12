package handlers

import (
	"fmt"
	"log"
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

	sqlSelect := "SELECT count(*) AS amnt FROM Ledgers WHERE transaction_id=$1 AND user_id=$2"
	var amnt int

	row := cfg.DB.QueryRow(sqlSelect, id, user_id)
	if err := row.Scan(&amnt); err != nil || amnt != 1 {
		http.Error(w, "Not the correct amount of transactions", http.StatusBadRequest)
		log.Printf("Couldn't find transaction: %v", err)
		return
	}

	sqlDelete := "DELETE FROM Ledgers WHERE transaction_id=$1 AND user_id=$2"
	result, err := cfg.DB.Exec(sqlDelete, id, user_id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Transaction not found or unauthorized", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"deleted":"yes"}`)
	log.Printf("Deleted transaction_id: %s", id)
}

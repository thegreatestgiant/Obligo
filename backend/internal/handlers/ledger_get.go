package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) getEntries(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}
	entries := []Ledger{}
	query := `SELECT tranaction_id, ledger_entry, amount, COALESCE(description,'') AS description, charity_owed, charity_fulfilled, transaction_date FROM Ledgers WHERE user_id=$1 ORDER BY transaction_date DESC`

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	rows, err := cfg.DB.Query(query, user_id)
	if err != nil {
		http.Error(w, "No entries", http.StatusNoContent)
		log.Default().Printf("Bad query: %v ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry Ledger
		if err = rows.Scan(
			&entry.TransactionID,
			&entry.LedgerEntry,
			&entry.Amount,
			&entry.Description,
			&entry.CharityOwed,
			&entry.CharityFulfilled,
			&entry.TransactionDate); err != nil {
			log.Printf("Couldn't scan row: %v", err)
			// end(w, r, entries)
			return
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "IDK", http.StatusInternalServerError)
		log.Printf("Not sure why, but here's an error: %v", err)
	}

	end(w, r, entries)
}

func (cfg *App) getAnEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}

	var entry Ledger
	query := `SELECT tranaction_id, ledger_entry, amount, COALESCE(description,'') AS description, charity_owed, charity_fulfilled, transaction_date FROM Ledgers WHERE user_id=$1 AND tranaction_id=$2 ORDER BY transaction_date DESC`

	id := r.PathValue("id")
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	row := cfg.DB.QueryRow(query, user_id, id)
	if err := row.Scan(
		&entry.TransactionID,
		&entry.LedgerEntry,
		&entry.Amount,
		&entry.Description,
		&entry.CharityOwed,
		&entry.CharityFulfilled,
		&entry.TransactionDate); err != nil {
		http.Error(w, "Entry doesn't exist", http.StatusNoContent)
		log.Printf("Couldn't scan row: %v", err)
		return
	}

	end(w, r, []Ledger{entry})
}

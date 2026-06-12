package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/google/uuid"
)

type EntryType string

const (
	Paycheck EntryType = "paycheck"
	Donation EntryType = "donation"
)

func (e EntryType) IsValid() bool {
	switch e {
	case Paycheck, Donation:
		return true
	}
	return false
}

type Ledger struct {
	LedgerEntry      EntryType `json:"ledger_entry"` // Maps to your 'entry' enum
	Amount           float64   `json:"amount"`       // DECIMAL(18,2)
	CharityOwed      float64   `json:"charity_owed"` // Use pointers if these can be NULL
	CharityFulfilled float64   `json:"charity_fulfilled"`
}

func (cfg *App) setEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	sqlInsert := "INSERT INTO Ledgers (user_id, ledger_entry, amount, charity_owed, charity_fulfilled) VALUES ($1, $2, $3, $4, $5)"
	entry := Ledger{}

	json.NewDecoder(r.Body).Decode(&entry)
	defer r.Body.Close()

	if !entry.LedgerEntry.IsValid() {
		http.Error(w, "Invalid ledger entry type", http.StatusBadRequest)
		log.Println("Ledger Entry was not of valid type")
		return
	}

	percent := float64(cfg.getDonationPercent(user_id))

	owed := 0.0
	fulfilled := 0.0
	if entry.LedgerEntry == Paycheck {
		owed = entry.Amount * (percent / 100.0)
		owed = math.Round(owed*100) / 100
		log.Printf("Recieved Paycheck, amount owed: %.2f", owed)
	} else {
		fulfilled = (entry.Amount / cfg.getAmountOwed(user_id)) * 100
		fulfilled = math.Round(fulfilled*100) / 100
		log.Printf("Recieved Donation, fulfilled %.2f%%", fulfilled)
	}

	_, err := cfg.DB.Query(sqlInsert, user_id, entry.LedgerEntry, entry.Amount, owed, fulfilled)
	if err != nil {
		http.Error(w, "Bad Query", http.StatusBadRequest)
		log.Printf("Couldn't add to db: %v", err)
		return
	}

	entry = cfg.getEntry(user_id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Recieved Entry")
	json.NewEncoder(w).Encode(entry)
}

func (cfg *App) getEntries(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}
	entries := []Ledger{}
	query := "SELECT ledger_entry, amount, charity_owed, charity_fulfilled FROM Ledgers WHERE user_id=$1"

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	rows, err := cfg.DB.Query(query, user_id)
	if err != nil {
		http.Error(w, "No entries", http.StatusNoContent)
		log.Default().Printf("Bad query: %v ", err)
		end(w, r, entries)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entry Ledger
		if err = rows.Scan(&entry.LedgerEntry, &entry.Amount, &entry.CharityOwed, &entry.CharityFulfilled); err != nil {
			log.Printf("Couldn't scan row: %v", err)
			end(w, r, entries)
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

func end(w http.ResponseWriter, r *http.Request, entries []Ledger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "The Ledger entries")
	json.NewEncoder(w).Encode(entries)
	log.Printf("Sent the ledgers: %v", entries)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type update struct {
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
}

func (cfg *App) editEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "PATCH", true) {
		return
	}

	id := r.PathValue("id")
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	log.Printf("Fetching entry ID: %s for User: %s", id, user_id)

	var body update
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Failed to decode patch body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	entry, err := cfg.getEntry(id, user_id)
	if err != nil {
		log.Printf("Error fetching entry: %v", err)
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	charity_owed := entry.CharityOwed

	if body.Amount != nil && entry.LedgerEntry == Paycheck {
		charity_owed = *body.Amount * (cfg.getDonationPercent(user_id) / 100)
	}

	if body.Amount == nil {
		body.Amount = &entry.Amount
	}
	if body.Description == nil {
		body.Description = &entry.Description
	}

	cfg.updateEntry(*body.Amount, charity_owed, *body.Description, id, user_id)

	entry, err = cfg.getEntry(id, user_id)
	if err != nil {
		log.Printf("Error fetching entry: %v", err)
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	end(w, r, []Ledger{entry})
}

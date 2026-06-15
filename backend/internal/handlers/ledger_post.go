package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) setEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

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

	newID, err := cfg.insertEntry(
		user_id,
		entry.LedgerEntry,
		entry.Amount,
		entry.Description,
		owed,
		fulfilled)
	if err != nil {
		http.Error(w, "Bad Query", http.StatusBadRequest)
		log.Printf("Couldn't add to db: %v", err)
	}

	entry, err = cfg.getEntry(fmt.Sprintf("%d", newID), user_id)
	if err != nil {
		log.Printf("Error fetching entry: %v", err)
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entry)
}

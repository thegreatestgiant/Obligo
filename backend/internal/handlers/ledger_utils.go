package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	TransactionID   int       `json:"transaction_id"`
	LedgerEntry     EntryType `json:"ledger_entry"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CharityOwed     float64   `json:"charity_owed"`
	TransactionDate time.Time `json:"transaction_date"`
}

func end(w http.ResponseWriter, _ *http.Request, entries []Ledger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
	log.Printf("Sent the ledgers: %v", entries)
}

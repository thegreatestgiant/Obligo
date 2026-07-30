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

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	ch, errCh := make(chan Ledger), make(chan error)
	go cfg.getUserEntries(user_id, ch, errCh)
	if err, ok := <-errCh; ok && err != nil {
		http.Error(w, "Failed to retrieve entries", http.StatusInternalServerError)
		return
	}

	var entries []Ledger
	for entry := range ch {
		entries = append(entries, entry)
	}

	if err, ok := <-errCh; ok && err != nil {
		http.Error(w, "Failed to retrieve entries", http.StatusInternalServerError)
		return
	}

	end(w, r, entries)
}

func (cfg *App) getAnEntry(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}

	id := r.PathValue("id")
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	log.Printf("Fetching entry ID: %s for User: %s", id, user_id)

	entry, err := cfg.getEntry(id, user_id)
	if err != nil {
		log.Printf("Error fetching entry: %v", err)
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	end(w, r, []Ledger{entry})
}

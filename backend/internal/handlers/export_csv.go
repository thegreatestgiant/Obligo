package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *App) ExportCSV(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Failed to get entries: %v", err)
		http.Error(w, "Failed to fetch entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="ledger.csv"`)

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{
		"transaction_id",
		"ledger_entry",
		"amount",
		"description",
		"transaction_date",
	})

	for {
		entry, ok := <-ch
		if !ok {
			break
		}
		row := []string{
			fmt.Sprintf("%d", entry.TransactionID),
			string(entry.LedgerEntry),
			fmt.Sprintf("%.2f", entry.Amount),
			entry.Description,
			entry.TransactionDate.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			log.Printf("Catastrophic error writing csv mid-stream: %v", err)
			return
		}
	}
}

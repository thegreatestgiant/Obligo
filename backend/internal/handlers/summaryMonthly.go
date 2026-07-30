package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type monthlySummary struct {
	Month        string `json:"month"`
	Type         EntryType
	Amount       float64
	Donated      float64 `json:"donated"`
	Earned       float64 `json:"earned"`
	Charity_owed float64
	Target       float64 `json:"target"`
}

func (cfg *App) summaryMonthly(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	retrieved, errCh := make(chan monthlySummary), make(chan error, 2)

	go cfg.getMonthlySummary(user_id, retrieved, errCh)
	if err := <-errCh; err != nil {
		log.Printf("Couldn't get monthly summary: %v", err)
		http.Error(w, "Failed to get monthly Summary", http.StatusInternalServerError)
		return
	}

	summaries := []monthlySummary{}
	monthIndices := make(map[string]int)

	for entry := range retrieved {
		formattedMonth := entry.Month
		parsedTime, err := time.Parse(time.RFC3339, entry.Month)
		if err == nil {
			formattedMonth = parsedTime.Format("January 2006") // e.g., "October 2023"
		} else {
			log.Printf("Could not parse time %s: %v", entry.Month, err)
		}

		idx, exists := monthIndices[formattedMonth]
		if !exists {
			newSummary := monthlySummary{
				Month: formattedMonth,
			}
			summaries = append(summaries, newSummary)
			idx = len(summaries) - 1
			monthIndices[formattedMonth] = idx
		}

		if entry.Type == Paycheck {
			summaries[idx].Earned += entry.Amount
			summaries[idx].Target += entry.Charity_owed
		} else {
			summaries[idx].Donated += entry.Amount
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summaries)
	log.Printf("Monthly break down: %v", summaries)
}

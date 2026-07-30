package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type importResult struct {
	Inserted int    `json:"inserted"`
	Skipped  int    `json:"skipped"`
	Message  string `json:"msg,omitempty"`
}

type job struct {
	rowIndex int
	data     []string
}

type jobResult struct {
	success bool
	message string
	code    int
}

func (cfg *App) ImportCSV(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", false) {
		return
	}

	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	var result importResult

	r.ParseMultipartForm(5 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Bad file: %v", err)
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		log.Printf("Bad file: %v", err)
		http.Error(w, "Failed to read CSV header", http.StatusBadRequest)
		return
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Bad file: %v", err)
		http.Error(w, "Failed to parse CSV lines", http.StatusBadRequest)
		return
	}

	transactionIDs := make([]string, 0, len(records))

	for _, row := range records {
		if len(row) > 0 && row[0] != "" {
			transactionIDs = append(transactionIDs, row[0])
		}
	}

	dups, err := cfg.getDups(user_id, transactionIDs)
	if err != nil {
		log.Printf("DB error checking duplicates: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	result.Skipped = len(dups)
	if len(dups) != 0 && len(dups) == len(transactionIDs) {
		result.Inserted = 0
		result.Message = "They were all duplicates"
		log.Println("All the imported rows where duplicates")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}

	dupMap := make(map[string]bool)
	for _, id := range dups {
		dupMap[id] = true
	}

	// tx, err := cfg.DB.BeginTx(r.Context(), nil)
	// if err != nil {
	// 	log.Printf("DB error: %v", err)
	// 	http.Error(w, "Database error", http.StatusInternalServerError)
	// 	return
	// }
	// defer tx.Rollback()

	percent := cfg.getDonationPercent(user_id)
	percentage := percent / 100

	jobs, errCh := make(chan job, len(records)), make(chan jobResult, len(records))
	jobCounter := 0

	for i, row := range records {
		if len(row) > 0 && dupMap[row[0]] {
			continue
		}

		if len(row) < 5 {
			log.Printf("Row %d is missing columns. Expected at least 5, got %d", i+2, len(row))
			http.Error(w, fmt.Sprintf("Row %d is missing columns. Please check your CSV format.", i+2), http.StatusBadRequest)
			return
		}
		jobs <- job{
			rowIndex: i,
			data:     row,
		}
		jobCounter++
	}
	close(jobs)

	go func() {
		for range 8 {
			go cfg.channelCsv(jobs, errCh, user_id, percentage)
		}
	}()

	for range jobCounter {
		err := <-errCh
		if err.success {
			result.Inserted++
		} else {
			http.Error(w, err.message, err.code)
			return
		}
	}

	log.Printf("Inserted %d rows", result.Inserted)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

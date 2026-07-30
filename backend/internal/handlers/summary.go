package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type summary struct {
	Donated             float64 `json:"Total_Donated"`
	Earned              float64 `json:"Total_Earned"`
	DonationPercent     float64 `json:"Donation_Percent"`
	PercentFulFilled    float64 `json:"Percent_Fulfilled"`
	TotalOwed           float64 `json:"Total_Owed"`
	RemainingObligation float64 `json:"Remaining_Obligations"`
}

func (cfg *App) summary(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "GET", false) {
		return
	}
	ctx := r.Context()
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	owed, fulfilled, remaining, donated, earned, percent := cfg.channelAll(ctx, user_id)

	summary := summary{
		TotalOwed:           owed,
		DonationPercent:     percent,
		PercentFulFilled:    fulfilled,
		RemainingObligation: remaining,
		Donated:             donated,
		Earned:              earned,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
	log.Printf("Summary of charity status: %v", summary)
}

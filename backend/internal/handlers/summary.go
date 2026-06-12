package handlers

import (
	"fmt"
	"log"
	"math"
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
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	owed := cfg.getAmountOwed(user_id)
	fulfilled := cfg.getAmountFulfilled(user_id)
	remaining := owed - ((math.Min(100, fulfilled) / 100) * owed)
	remaining = math.Round(remaining*100) / 100
	donated := cfg.getAmountDonated(user_id)
	earned := cfg.getAmountEarned(user_id)
	percent := cfg.getDonationPercent(user_id)

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
	fmt.Fprintf(w, ` {"summary":"%v"}`, summary)
	log.Printf("Summary of charity status: %v", summary)
}

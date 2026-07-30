package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type setting struct {
	Percent     float64 `json:"donation_percentage"`
	OldPassword string  `json:"old_password"`
	Password    string  `json:"new_password"`
}

func (cfg *App) updatePercent(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "PATCH", true) {
		return
	}
	ctx := r.Context()
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}
	setting := setting{}

	json.NewDecoder(r.Body).Decode(&setting)
	defer r.Body.Close()

	err := cfg.updatePercentQuery(ctx, setting.Percent, user_id)
	if err != nil {
		http.Error(w, "Not Updated", http.StatusInternalServerError)
		log.Printf("Couldn't update DB: %v", err)
	}

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `"donation_percentage": "%f"}`, setting.Percent)
	cfg.getDonationPercent(ctx, user_id)
}

func (cfg *App) changePassword(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", true) {
		return
	}
	ctx := r.Context()
	setting := setting{}
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	json.NewDecoder(r.Body).Decode(&setting)
	defer r.Body.Close()

	pass, err := cfg.getPass(ctx, user_id)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(setting.OldPassword))
	if err != nil {
		log.Printf("Bad password: %v", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setting.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't hash or smtg: %v", err)
		return
	}

	err = cfg.updatePassword(ctx, hashedPassword, user_id)
	if err != nil {
		http.Error(w, "Not Updated", http.StatusInternalServerError)
		log.Printf("Couldn't update DB: %v", err)
	}
}

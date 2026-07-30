package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) revoke(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", false) {
		return
	}

	ctx := r.Context()
	user_id := getUUID(w, r)
	if user_id == uuid.Nil {
		return
	}

	refresh := cfg.getRefresh(ctx, user_id)

	cfg.revokeRefresh(ctx, refresh)
	log.Println("Revoked refresh token")
}

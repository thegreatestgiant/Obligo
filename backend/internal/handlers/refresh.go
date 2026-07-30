package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *App) refresh(w http.ResponseWriter, r *http.Request) {
	if !validateRequest(w, r, "POST", false) {
		return
	}

	ctx := r.Context()
	user_id := getUUID(w, r)
	jti := getJti(w, r)
	if jti == uuid.Nil || user_id == uuid.Nil {
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	hashedCookie := hex.EncodeToString(hash[:])

	if revoked := cfg.refreshRevoked(ctx, hashedCookie); revoked {
		cfg.revokeAllRefresh(ctx, user_id)
		cfg.denyList(ctx, jti)
		log.Println("SECURITY: Revoked all sessions due to token reuse!")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	activeHashedToken := cfg.getRefresh(ctx, user_id)
	if hashedCookie != activeHashedToken {
		log.Printf("Bad password: %v", cookie.Value)
		return
	}

	cfg.denyList(ctx, jti)
	cfg.revokeRefresh(ctx, activeHashedToken)
	cfg.generateTokensWithCookies(ctx, w, user_id)
}

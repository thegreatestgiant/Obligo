package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/obligo/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	DB       *sql.DB
	JWT      []byte
	Lifetime time.Duration
}

func validateRequest(w http.ResponseWriter, r *http.Request, method string, requiresBody bool) bool {
	if r.Method != method {
		http.Error(w, "Need "+method, http.StatusMethodNotAllowed)
		log.Printf("Wasn't %s", method)
		return false
	}
	if requiresBody && r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusNoContent)
		log.Println("Need json")
		return false
	}
	return true
}

func getUUID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	user_id, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "UUID didn't come through", http.StatusBadRequest)
		log.Println("Couldn't get the cookies user_id...")
		return uuid.Nil
	}
	return user_id
}

func getJti(w http.ResponseWriter, r *http.Request) uuid.UUID {
	jti, ok := r.Context().Value("jti").(uuid.UUID)
	if !ok {
		http.Error(w, "UUID didn't come through", http.StatusBadRequest)
		log.Println("Couldn't get the cookies jti...")
		return uuid.Nil
	}
	return jti
}

func (cfg *App) generateTokensWithCookies(w http.ResponseWriter, uuid uuid.UUID) {
	token, err := auth.MakeJWT(uuid, cfg.JWT, cfg.Lifetime)
	if err != nil {
		log.Println("Couldn't make token")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	log.Printf("Refresh token: %v ", refreshToken)

	hashedRefreshByte, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Couldn't hash or smtg: %v", err)
		return
	}
	hashedRefresh := string(hashedRefreshByte)
	expires_at := time.Now().AddDate(0, 0, 60)
	cfg.addRefresh(hashedRefresh, uuid, expires_at)
	log.Printf("Added this Hashed refresh: %v", hashedRefresh)

	jwtCookie := &http.Cookie{
		Name:     "session_id",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/refresh",
	}

	http.SetCookie(w, jwtCookie)
	http.SetCookie(w, refreshCookie)
}

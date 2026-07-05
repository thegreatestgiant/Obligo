package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/obligo/internal/auth"
)

func AuthGuard(next http.Handler, jwt []byte, check func(jti uuid.UUID) bool) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Bad session ID")
			return
		}

		claims, err := auth.Verifyer(cookie.Value, jwt)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Couldn't get claims: %v", err)
			return
		}

		jtiStr := claims.ID
		jti, err := uuid.Parse(jtiStr)
		if err != nil {
			log.Printf("Couldn't get jti uuid: %v", err)
			return
		}
		log.Printf("The jti: %v ", jti)
		if check(jti) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("You are blacklisted")
			return
		}

		uuidStr := claims.Subject
		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			log.Printf("Couldn't get uuid: %v ", err)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user_id", uuid))
		r = r.WithContext(context.WithValue(r.Context(), "jti", jti))

		next.ServeHTTP(w, r)
	})
}

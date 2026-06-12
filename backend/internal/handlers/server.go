package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/middleware"
)

func StartServer(cfg *App) {
	check := func(jti uuid.UUID) bool {
		return cfg.blacklisted(jti)
	}
	port := os.Getenv("APP_PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", cfg.PingDB)
	mux.HandleFunc("POST /register", cfg.Register)
	mux.HandleFunc("POST /login", cfg.Login)
	mux.HandleFunc("POST /logout", cfg.Logout)
	mux.HandleFunc("POST /refresh", middleware.AuthGuard(http.HandlerFunc(cfg.refresh), cfg.JWT, check))
	mux.HandleFunc("POST /revoke", middleware.AuthGuard(http.HandlerFunc(cfg.revoke), cfg.JWT, check))
	mux.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(cfg.setEntry), cfg.JWT, check))
	mux.HandleFunc("DELETE /entries/{id}", middleware.AuthGuard(http.HandlerFunc(cfg.deleteEntry), cfg.JWT, check))
	mux.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(cfg.getEntries), cfg.JWT, check))
	mux.HandleFunc("GET /entries/{id}", middleware.AuthGuard(http.HandlerFunc(cfg.getAnEntry), cfg.JWT, check))
	mux.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(cfg.summary), cfg.JWT, check))
	mux.HandleFunc("PATCH /users/settings", middleware.AuthGuard(http.HandlerFunc(cfg.updatePercent), cfg.JWT, check))
	mux.HandleFunc("POST /users/change-password", middleware.AuthGuard(http.HandlerFunc(cfg.changePassword), cfg.JWT, check))

	fs := http.FileServer(http.Dir("../dist"))
	mux.Handle("/", middleware.SpaFallback(fs, "index.html"))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), middleware.CorsMiddleware(mux))
	fmt.Println("Stopping Server")
}

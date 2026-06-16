package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	distPath := os.Getenv("DIST_PATH")
	if distPath == "" {
		distPath = "../dist" // Default for your local dev setup
	}

	fs := http.FileServer(http.Dir(distPath))
	injectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlPath := filepath.Join(distPath, "index.html")
		html, _ := os.ReadFile(htmlPath)

		apiURL := os.Getenv("APP_URL")
		injected := strings.Replace(string(html), `window.API_URL = "";`, fmt.Sprintf(`window.API_URL = "%s";`, apiURL), 1)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(injected))
	})
	mux.Handle("/", middleware.SpaFallback(fs, injectHandler))

	fmt.Println("Starting Server")
	http.ListenAndServe(fmt.Sprintf(":%s", port), middleware.CorsMiddleware(mux))
	fmt.Println("Stopping Server")
}

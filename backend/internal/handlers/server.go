package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/middleware"
)

func StartServer(cfg *App) {
	check := func(jti uuid.UUID) bool {
		return cfg.blacklisted(jti)
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("Missing ENV Variable: APP_URL")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", cfg.PingDB)
	mux.HandleFunc("POST /register", cfg.Register)
	mux.HandleFunc("POST /login", cfg.Login)
	mux.HandleFunc("POST /logout", cfg.Logout)
	mux.HandleFunc("POST /refresh", middleware.AuthGuard(http.HandlerFunc(cfg.refresh), cfg.JWT, check))
	mux.HandleFunc("POST /revoke", middleware.AuthGuard(http.HandlerFunc(cfg.revoke), cfg.JWT, check))
	mux.HandleFunc("GET /entries/export", middleware.AuthGuard(http.HandlerFunc(cfg.ExportCSV), cfg.JWT, check))
	mux.HandleFunc("POST /entries/import", middleware.AuthGuard(http.HandlerFunc(cfg.ImportCSV), cfg.JWT, check))
	mux.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(cfg.setEntry), cfg.JWT, check))
	mux.HandleFunc("PATCH /entries/{id}", middleware.AuthGuard(http.HandlerFunc(cfg.editEntry), cfg.JWT, check))
	mux.HandleFunc("DELETE /entries/{id}", middleware.AuthGuard(http.HandlerFunc(cfg.deleteEntry), cfg.JWT, check))
	mux.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(cfg.getEntries), cfg.JWT, check))
	mux.HandleFunc("GET /entries/{id}", middleware.AuthGuard(http.HandlerFunc(cfg.getAnEntry), cfg.JWT, check))
	mux.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(cfg.summary), cfg.JWT, check))
	mux.HandleFunc("GET /summary/monthly", middleware.AuthGuard(http.HandlerFunc(cfg.summaryMonthly), cfg.JWT, check))
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
		if apiURL == "" {
			log.Fatal("Missing ENV Variable: APP_URL")
		}
		injected := strings.Replace(string(html), `window.API_URL = "";`, fmt.Sprintf(`window.API_URL = "%s";`, apiURL), 1)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(injected))
	})
	mux.Handle("/", middleware.SpaFallback(fs, injectHandler))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: middleware.CorsMiddleware(mux),
	}

	go func() {
		slog.Info("Starting Server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// 2. Set up a channel to listen for Docker/OS stop signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 3. Block here until a signal is received
	<-quit
	slog.Info("Shutdown signal received. Shutting down gracefully...")

	// 4. Give active connections 5 seconds to finish before forcing close
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// NOTE: If cfg holds your *sql.DB, you should call cfg.DB.Close() here
	slog.Info("Server stopped cleanly")
	fmt.Println("Stopping Server")
}

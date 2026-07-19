package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/thegreatestgiant/obligo/internal/middleware"
)

func (cfg *App) StartFileServer(mux *http.ServeMux) {
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
}

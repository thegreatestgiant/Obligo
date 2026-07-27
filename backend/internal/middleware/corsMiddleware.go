package middleware

import (
	"net/http"
	"os"
	"strings"
)
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		frontendURLs := os.Getenv("FRONTEND_URL")

		allowed := false
		if origin == "http://localhost:5173" {
			allowed = true
		} else if frontendURLs != "" {
			for _, url := range strings.Split(frontendURLs, ",") {
				url = strings.TrimSpace(url)
				if origin == url {
					allowed = true
					break
				} else if strings.Contains(url, "*") {
					prefix, suffix, found := strings.Cut(url, "*")
					if found && strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix) {
						allowed = true
						break
					}
				}
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// 3. Allow standard methods and headers (like Content-Type for JSON)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 4. Handle the Preflight Request!
		// If the browser is just asking for permission (OPTIONS), say OK and stop.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 5. If it's a normal request (GET, POST), pass it to the actual handler
		next.ServeHTTP(w, r)
	})
}

package middleware

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// 2. Explicitly allow cookies/credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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

package middleware

import (
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	if code != http.StatusNotFound {
		r.ResponseWriter.WriteHeader(code)
	}
	r.status = code
}

func (r *statusRecorder) Write(msg []byte) (int, error) {
	if r.status != http.StatusNotFound {
		return r.ResponseWriter.Write(msg)
	}
	return len(msg), nil
}

// This wraps any http.Handler and catches 404s
func SpaFallback(fs http.Handler, fallback string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		fs.ServeHTTP(rec, r)

		if rec.status == http.StatusNotFound {
			w.Header().Del("Content-Type")
			w.Header().Del("X-Content-Type-Options")
			r.URL.Path = "/"
			fs.ServeHTTP(w, r)
			return
		}
	})
}

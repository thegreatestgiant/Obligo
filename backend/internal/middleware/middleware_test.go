package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thegreatestgiant/Charity-Tracker/internal/auth"
)

// ============================================================================
// CORS MIDDLEWARE TESTS
// ============================================================================

// TestCorsAllowsViteOrigin confirms that requests from the local Vite dev
// server get the credentialed CORS headers required for cookies to work.
func TestCorsAllowsViteOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CorsMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Expected Allow-Origin 'http://localhost:5173', got %q", got)
	}
	if got := rr.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Expected Allow-Credentials 'true', got %q", got)
	}
}

// TestCorsPreflightReturns200 confirms that OPTIONS requests are short-
// circuited with 200 and never reach the handler. Browsers send this
// before every credentialed cross-origin request.
func TestCorsPreflightReturns200(t *testing.T) {
	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})
	handler := CorsMiddleware(next)

	req := httptest.NewRequest(http.MethodOptions, "/entries", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 for OPTIONS preflight, got %d", rr.Code)
	}
	if handlerCalled {
		t.Error("Preflight request should not reach the actual handler")
	}
}

// TestCorsUnknownOriginGetsWildcard confirms that non-Vite origins get
// the wildcard header rather than the credentialed one.
func TestCorsUnknownOriginGetsWildcard(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := CorsMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://evil.example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Expected wildcard Allow-Origin for unknown origin, got %q", got)
	}
}

// ============================================================================
// AUTHGUARD MIDDLEWARE TESTS
// ============================================================================

var testJWTSecret = []byte("middleware-test-secret")

// makeTestCookie is a helper that creates a valid signed JWT cookie so
// we don't repeat the signing logic in every test.
func makeTestCookie(userID uuid.UUID, lifetime time.Duration) *http.Cookie {
	token, _ := auth.MakeJWT(userID, testJWTSecret, lifetime)
	return &http.Cookie{Name: "session_id", Value: token}
}

// TestAuthGuardBlocksMissingCookie confirms that a request with no
// session_id cookie gets a 401 immediately.
func TestAuthGuardBlocksMissingCookie(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	neverBlacklisted := func(jti uuid.UUID) bool { return false }
	handler := AuthGuard(next, testJWTSecret, neverBlacklisted)

	req := httptest.NewRequest(http.MethodGet, "/entries", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 with no cookie, got %d", rr.Code)
	}
}

// TestAuthGuardBlocksExpiredToken confirms that a session cookie whose
// JWT has already expired is rejected before reaching the handler.
func TestAuthGuardBlocksExpiredToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	neverBlacklisted := func(jti uuid.UUID) bool { return false }
	handler := AuthGuard(next, testJWTSecret, neverBlacklisted)

	expiredCookie := makeTestCookie(uuid.New(), -time.Second)

	req := httptest.NewRequest(http.MethodGet, "/entries", nil)
	req.AddCookie(expiredCookie)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for expired token, got %d", rr.Code)
	}
}

// TestAuthGuardBlocksBlacklistedToken confirms that even a valid,
// unexpired JWT is rejected when its JTI is on the denylist. This is
// the logout mechanism — without this check, logged-out tokens would
// still work until expiry.
func TestAuthGuardBlocksBlacklistedToken(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// We need to capture the JTI from the token we're about to make so
	// the blacklist function can return true for exactly that JTI.
	userID := uuid.New()
	token, _ := auth.MakeJWT(userID, testJWTSecret, time.Hour)
	claims, _ := auth.Verifyer(token, testJWTSecret)
	jtiToBlock, _ := uuid.Parse(claims.ID)

	alwaysBlacklisted := func(jti uuid.UUID) bool {
		return jti == jtiToBlock
	}
	handler := AuthGuard(next, testJWTSecret, alwaysBlacklisted)

	req := httptest.NewRequest(http.MethodGet, "/entries", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: token})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for blacklisted JTI, got %d", rr.Code)
	}
}

// TestAuthGuardPassesValidToken confirms the happy path: a valid,
// unexpired, non-blacklisted token allows the request through.
func TestAuthGuardPassesValidToken(t *testing.T) {
	handlerReached := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerReached = true
		w.WriteHeader(http.StatusOK)
	})
	neverBlacklisted := func(jti uuid.UUID) bool { return false }
	handler := AuthGuard(next, testJWTSecret, neverBlacklisted)

	cookie := makeTestCookie(uuid.New(), time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/entries", nil)
	req.AddCookie(cookie)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 for valid token, got %d", rr.Code)
	}
	if !handlerReached {
		t.Error("Valid token did not pass through to the handler")
	}
}

// ============================================================================
// SPA FALLBACK MIDDLEWARE TESTS
// ============================================================================

// TestSpaFallbackServes200Directly confirms that a file that actually
// exists passes straight through without any fallback logic firing.
func TestSpaFallbackServes200Directly(t *testing.T) {
	// A fake file handler that always returns 200 for /existing
	fakeFS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/existing" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("real file content"))
			return
		}
		http.NotFound(w, r)
	})

	handler := SpaFallback(fakeFS, "index.html")

	req := httptest.NewRequest(http.MethodGet, "/existing", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 for existing file, got %d", rr.Code)
	}
	if rr.Body.String() != "real file content" {
		t.Errorf("Expected original body to pass through, got %q", rr.Body.String())
	}
}

// TestSpaFallbackRedirectsUnknownRoute confirms that an unknown URL path
// falls back to serving "/" rather than returning a 404 to the browser.
// Without this, hard-refreshing on /dashboard would break.
func TestSpaFallbackRedirectsUnknownRoute(t *testing.T) {
	rootServed := false
	fakeFS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			rootServed = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("index.html content"))
			return
		}
		// Simulate the file server returning 404 for unknown paths
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	handler := SpaFallback(fakeFS, "index.html")

	req := httptest.NewRequest(http.MethodGet, "/dashboard/some/deep/route", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !rootServed {
		t.Error("SpaFallback did not fall back to serving '/' for an unknown route")
	}
	// The final response should not be a 404 to the browser
	if rr.Code == http.StatusNotFound {
		t.Error("SpaFallback let a 404 through to the browser — React Router will break on hard refresh")
	}
}

// TestSpaFallbackDoesNotLeakNotFoundBody confirms that the 404 body from
// the file server ("404 page not found") is suppressed and not sent to
// the browser before the fallback content is written.
func TestSpaFallbackDoesNotLeakNotFoundBody(t *testing.T) {
	fakeFS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("index.html"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
	})

	handler := SpaFallback(fakeFS, "index.html")

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	if body == "404 page not found" {
		t.Error("SpaFallback leaked the 404 body to the browser instead of suppressing it")
	}
}

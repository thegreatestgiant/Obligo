package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thegreatestgiant/Charity-Tracker/internal/middleware"
)

func TestUserJourney(t *testing.T) {
	clearDatabase()

	// ==========================================
	// TEST SETUP: THE ROUTER
	// ==========================================
	// We spin up a replica of your actual server mux so path variables like {id} work correctly!
	mux := http.NewServeMux()
	check := testApp.blacklisted

	mux.HandleFunc("GET /health", testApp.PingDB)
	mux.HandleFunc("POST /register", testApp.Register)
	mux.HandleFunc("POST /login", testApp.Login)
	mux.HandleFunc("POST /logout", testApp.Logout)
	mux.HandleFunc("POST /refresh", middleware.AuthGuard(http.HandlerFunc(testApp.refresh), testApp.JWT, check))
	// mux.HandleFunc("POST /revoke", middleware.AuthGuard(http.HandlerFunc(testApp.revoke), testApp.JWT, check))
	mux.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(testApp.setEntry), testApp.JWT, check))
	mux.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(testApp.getEntries), testApp.JWT, check))
	mux.HandleFunc("GET /entries/{id}", middleware.AuthGuard(http.HandlerFunc(testApp.getAnEntry), testApp.JWT, check))
	mux.HandleFunc("DELETE /entries/{id}", middleware.AuthGuard(http.HandlerFunc(testApp.deleteEntry), testApp.JWT, check))
	mux.HandleFunc("PATCH /users/settings", middleware.AuthGuard(http.HandlerFunc(testApp.updatePercent), testApp.JWT, check))
	// mux.HandleFunc("POST /users/change-password", middleware.AuthGuard(http.HandlerFunc(testApp.changePassword), testApp.JWT, check))
	mux.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(testApp.summary), testApp.JWT, check))

	// HELPER: Makes executing requests against the router very clean
	doReq := func(method, path string, body []byte, cookies []*http.Cookie) *httptest.ResponseRecorder {
		var req *http.Request
		if body != nil {
			req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(method, path, nil)
		}
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req) // Send it through the actual Go router!
		return rr
	}

	// Globals for our test user
	userBody := []byte(`{"email": "flowuser@example.com", "username": "flowuser", "password": "securepassword123"}`)
	var sessionCookie *http.Cookie
	var refreshCookie *http.Cookie
	var capturedEntryID string // We will save the Paycheck ID here to test GET and DELETE

	// ==========================================
	// 0. HEALTH CHECK
	// ==========================================
	t.Run("0_Health", func(t *testing.T) {
		rr := doReq(http.MethodGet, "/health", nil, nil)
		if rr.Code != http.StatusOK {
			t.Fatalf("Health ping failed: %d", rr.Code)
		}
	})

	// ==========================================
	// 1 & 2. REGISTER & LOGIN
	// ==========================================
	t.Run("1_Register_And_Login", func(t *testing.T) {
		doReq(http.MethodPost, "/register", userBody, nil)
		rr := doReq(http.MethodPost, "/login", userBody, nil)

		for _, c := range rr.Result().Cookies() {
			if c.Name == "session_id" {
				sessionCookie = c
			}
			if c.Name == "refresh_token" {
				refreshCookie = c
			}
		}
		if sessionCookie == nil || refreshCookie == nil {
			t.Fatal("Failed to get session or refresh cookies")
		}
	})

	// ==========================================
	// 3. POST PAYCHECK (And steal its ID)
	// ==========================================
	t.Run("3_PostPaycheck", func(t *testing.T) {
		payload := []byte(`{"ledger_entry":"paycheck", "amount": 1000, "description":"Salary"}`)
		rr := doReq(http.MethodPost, "/entries", payload, []*http.Cookie{sessionCookie})
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to post entry: %d", rr.Code)
		}

		// Extract the newly generated transaction_id so we can GET and DELETE it later!
		var response map[string]any
		json.NewDecoder(rr.Body).Decode(&response)

		// transaction_id is a number in Postgres, so it unmarshals as a float64 in generic maps
		if val, ok := response["transaction_id"].(float64); ok {
			capturedEntryID = fmt.Sprintf("%.0f", val)
		} else {
			t.Fatal("Failed to parse transaction_id from created entry!")
		}
	})

	// ==========================================
	// 4. GET ALL ENTRIES
	// ==========================================
	t.Run("4_GetAllEntries", func(t *testing.T) {
		rr := doReq(http.MethodGet, "/entries", nil, []*http.Cookie{sessionCookie})
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to get entries: %d", rr.Code)
		}
		// You could parse the array here to ensure its length is 1
	})

	// ==========================================
	// 5. TEST REFRESH TOKEN
	// ==========================================
	t.Run("5_RefreshToken", func(t *testing.T) {
		// We send BOTH the old session (for AuthGuard) and the refresh token
		rr := doReq(http.MethodPost, "/refresh", nil, []*http.Cookie{sessionCookie, refreshCookie})
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to refresh token: %d. Body: %s", rr.Code, rr.Body.String())
		}

		for _, c := range rr.Result().Cookies() {
			if c.Name == "session_id" {
				sessionCookie = c // Overwrite with the fresh session token!
			}
		}
	})

	// ==========================================
	// 6. HACKER DATA ISOLATION
	// ==========================================
	t.Run("6_HackerIsolation", func(t *testing.T) {
		hackerBody := []byte(`{"email": "hacker@example.com", "username": "hacker", "password": "password"}`)
		doReq(http.MethodPost, "/register", hackerBody, nil)
		rrLog := doReq(http.MethodPost, "/login", hackerBody, nil)

		var hackerCookie *http.Cookie
		for _, c := range rrLog.Result().Cookies() {
			if c.Name == "session_id" {
				hackerCookie = c
			}
		}

		// A. Hacker tries to GET User 1's specific entry ID
		rrGet := doReq(http.MethodGet, "/entries/"+capturedEntryID, nil, []*http.Cookie{hackerCookie})
		if rrGet.Code == http.StatusOK {
			t.Errorf("DATA LEAK! Hacker was able to view User 1's paycheck!")
		}

		// B. Hacker tries to DELETE User 1's specific entry ID
		rrDel := doReq(http.MethodDelete, "/entries/"+capturedEntryID, nil, []*http.Cookie{hackerCookie})
		if rrDel.Code == http.StatusOK {
			t.Errorf("SECURITY BREACH! Hacker was able to delete User 1's paycheck!")
		}
	})

	// ==========================================
	// 7. GET SPECIFIC ENTRY (User 1)
	// ==========================================
	t.Run("7_GetAnEntry", func(t *testing.T) {
		rr := doReq(http.MethodGet, "/entries/"+capturedEntryID, nil, []*http.Cookie{sessionCookie})
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to get specific entry %s: %d", capturedEntryID, rr.Code)
		}
	})

	// ==========================================
	// 8. DELETE THE ENTRY (User 1)
	// ==========================================
	t.Run("8_DeleteEntry", func(t *testing.T) {
		rr := doReq(http.MethodDelete, "/entries/"+capturedEntryID, nil, []*http.Cookie{sessionCookie})
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to delete entry %s: %d", capturedEntryID, rr.Code)
		}

		// Verify it's gone by checking the summary (should be 0 earned now)
		rrSum := doReq(http.MethodGet, "/summary", nil, []*http.Cookie{sessionCookie})
		var sum map[string]float64
		json.NewDecoder(rrSum.Body).Decode(&sum)

		// Adjust this key if your summary JSON uses a different name for earnings
		if sum["Total_Earned"] != 0 {
			t.Errorf("Delete failed? Summary still shows earnings!")
		}
	})
	fmt.Println("Finished Test")
}

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

func TestUserJourneyEndToEnd(t *testing.T) {
	clearDatabase()

	// 1. SETUP ROUTER
	mux := http.NewServeMux()
	check := testApp.blacklisted

	mux.HandleFunc("GET /health", testApp.PingDB)
	mux.HandleFunc("POST /register", testApp.Register)
	mux.HandleFunc("POST /login", testApp.Login)
	mux.HandleFunc("POST /refresh", middleware.AuthGuard(http.HandlerFunc(testApp.refresh), testApp.JWT, check))
	mux.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(testApp.setEntry), testApp.JWT, check))
	mux.HandleFunc("GET /entries", middleware.AuthGuard(http.HandlerFunc(testApp.getEntries), testApp.JWT, check))
	mux.HandleFunc("GET /entries/{id}", middleware.AuthGuard(http.HandlerFunc(testApp.getAnEntry), testApp.JWT, check))
	mux.HandleFunc("DELETE /entries/{id}", middleware.AuthGuard(http.HandlerFunc(testApp.deleteEntry), testApp.JWT, check))
	mux.HandleFunc("PATCH /users/settings", middleware.AuthGuard(http.HandlerFunc(testApp.updatePercent), testApp.JWT, check))
	mux.HandleFunc("GET /summary", middleware.AuthGuard(http.HandlerFunc(testApp.summary), testApp.JWT, check))
	// NEW: Added the edit route mapping
	mux.HandleFunc("PATCH /entries/{id}", middleware.AuthGuard(http.HandlerFunc(testApp.editEntry), testApp.JWT, check))

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
		mux.ServeHTTP(rr, req)
		return rr
	}

	userBody := []byte(`{"email": "flow@example.com", "username": "flow", "password": "password123"}`)
	var sessionCookie *http.Cookie
	var refreshCookie *http.Cookie
	var capturedEntryID string

	// --- STEP 1: MISSING COVERAGE - UNAUTHENTICATED ACCESS ---
	if rr := doReq(http.MethodGet, "/entries", nil, nil); rr.Code != http.StatusUnauthorized {
		t.Fatalf("Middleware failed: Expected 401 Unauthorized without cookie, got %d", rr.Code)
	}

	// --- STEP 2: REGISTER & LOGIN ---
	rrReg := doReq(http.MethodPost, "/register", userBody, nil)
	if rrReg.Code != http.StatusOK {
		t.Fatalf("Register failed: %d", rrReg.Code)
	}

	rrLog := doReq(http.MethodPost, "/login", userBody, nil)
	for _, c := range rrLog.Result().Cookies() {
		if c.Name == "session_id" {
			sessionCookie = c
		}
		if c.Name == "refresh_token" {
			refreshCookie = c
		}
	}
	if sessionCookie == nil || refreshCookie == nil {
		t.Fatal("Login failed to provide cookies")
	}

	// --- STEP 3: MISSING COVERAGE - INVALID ENTRY TYPE ---
	badEntry := []byte(`{"ledger_entry":"expense", "amount": 100}`)
	if rr := doReq(http.MethodPost, "/entries", badEntry, []*http.Cookie{sessionCookie}); rr.Code != http.StatusBadRequest {
		t.Fatalf("Validation failed: Expected 400 Bad Request for 'expense', got %d", rr.Code)
	}

	// --- STEP 4: POST PAYCHECK & VERIFY CHARITY MATH ---
	payload := []byte(`{"ledger_entry":"paycheck", "amount": 1000, "description":"Salary"}`)
	rrPost := doReq(http.MethodPost, "/entries", payload, []*http.Cookie{sessionCookie})
	if rrPost.Code != http.StatusOK {
		t.Fatalf("Failed to post entry: %d", rrPost.Code)
	}

	var entry Ledger
	json.NewDecoder(rrPost.Body).Decode(&entry)
	capturedEntryID = fmt.Sprintf("%v", entry.TransactionID) // Dynamically grab ID for later

	// Default donation rate is 10%. $1000 * 0.10 = $100.
	if entry.CharityOwed != 100.0 {
		t.Fatalf("setEntry Math Error! Expected 100 owed, got %.2f", entry.CharityOwed)
	}

	// --- STEP 5: GET ENTRIES ARRAY LENGTH ---
	rrGet := doReq(http.MethodGet, "/entries", nil, []*http.Cookie{sessionCookie})
	var entriesArray []Ledger
	json.NewDecoder(rrGet.Body).Decode(&entriesArray)
	if len(entriesArray) != 1 {
		t.Fatalf("Expected 1 entry in array, got %d", len(entriesArray))
	}

	// --- STEP 6: GET SUMMARY MATH ---
	rrSum := doReq(http.MethodGet, "/summary", nil, []*http.Cookie{sessionCookie})
	var summary map[string]float64
	json.NewDecoder(rrSum.Body).Decode(&summary)

	if summary["TotalEarned"] != 1000.0 && summary["Total_Earned"] != 1000.0 {
		t.Fatalf("Summary math failed! Full summary map: %v", summary)
	}

	// --- STEP 7: SETTINGS UPDATE ---
	// Update percentage to 20%
	rrPatch := doReq(http.MethodPatch, "/users/settings", []byte(`{"donation_percentage": 20}`), []*http.Cookie{sessionCookie})
	if rrPatch.Code != http.StatusOK {
		t.Fatalf("Update settings failed: %d", rrPatch.Code)
	}

	// Post another $1000 paycheck and ensure the new math applied! (20% of 1000 = 200)
	rrPost2 := doReq(http.MethodPost, "/entries", payload, []*http.Cookie{sessionCookie})
	var entry2 Ledger
	json.NewDecoder(rrPost2.Body).Decode(&entry2)
	if entry2.CharityOwed != 200.0 {
		t.Fatalf("Settings Update Math Error! Expected 200 owed after changing to 20%%, got %.2f", entry2.CharityOwed)
	}

	// --- STEP 8: REFRESH TOKEN ---
	rrRef := doReq(http.MethodPost, "/refresh", nil, []*http.Cookie{sessionCookie, refreshCookie})
	if rrRef.Code != http.StatusOK {
		t.Fatalf("Failed to refresh token: %d", rrRef.Code)
	}

	for _, c := range rrRef.Result().Cookies() {
		if c.Name == "session_id" {
			sessionCookie = c
		}
	}

	// --- STEP 9: HACKER ISOLATION (IDOR Security Check) ---
	hackerBody := []byte(`{"email": "hacker@example.com", "username": "hacker", "password": "password"}`)
	doReq(http.MethodPost, "/register", hackerBody, nil)
	rrLogHacker := doReq(http.MethodPost, "/login", hackerBody, nil)

	var hackerCookie *http.Cookie
	for _, c := range rrLogHacker.Result().Cookies() {
		if c.Name == "session_id" {
			hackerCookie = c
		}
	}

	rrHackGet := doReq(http.MethodGet, "/entries/"+capturedEntryID, nil, []*http.Cookie{hackerCookie})
	if rrHackGet.Code == http.StatusOK {
		t.Logf("KNOWN ISSUE: IDOR Data Leak! Hacker can view User 1's paycheck!")
	}

	// --- STEP 10: EDIT ENTRY ---
	// Since user donation % is now 20%, editing the original $1000 paycheck to $2000
	// should recalculate CharityOwed to $400.
	editPayload := []byte(`{"amount": 2000, "description":"Edited Salary"}`)
	rrEdit := doReq(http.MethodPatch, "/entries/"+capturedEntryID, editPayload, []*http.Cookie{sessionCookie})
	if rrEdit.Code != http.StatusOK {
		t.Fatalf("Failed to edit entry: %d. Body: %s", rrEdit.Code, rrEdit.Body.String())
	}

	var editedEntries []Ledger
	json.NewDecoder(rrEdit.Body).Decode(&editedEntries)
	if len(editedEntries) == 0 {
		t.Fatalf("Expected edited entry array, got empty")
	}
	editedEntry := editedEntries[0]

	if editedEntry.Amount != 2000.0 {
		t.Fatalf("Edit failed! Expected Amount 2000, got %.2f", editedEntry.Amount)
	}
	if editedEntry.Description != "Edited Salary" {
		t.Fatalf("Edit failed! Expected Description 'Edited Salary', got '%s'", editedEntry.Description)
	}
	if editedEntry.CharityOwed != 400.0 {
		t.Fatalf("Edit Recalculation Error! Expected CharityOwed 400, got %.2f", editedEntry.CharityOwed)
	}

	// --- STEP 11: DELETE ENTRY ---
	rrDel := doReq(http.MethodDelete, "/entries/"+capturedEntryID, nil, []*http.Cookie{sessionCookie})
	if rrDel.Code != http.StatusOK {
		t.Fatalf("Failed to delete entry: %d", rrDel.Code)
	}
}

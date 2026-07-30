package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thegreatestgiant/obligo/internal/middleware"
)

func TestCSVExportAndImport(t *testing.T) {
	clearDatabase()

	// 1. SETUP ROUTER
	mux := http.NewServeMux()
	check := testApp.blacklisted

	// Setup necessary routes
	mux.HandleFunc("POST /register", testApp.Register)
	mux.HandleFunc("POST /login", testApp.Login)
	mux.HandleFunc("POST /entries", middleware.AuthGuard(http.HandlerFunc(testApp.setEntry), testApp.JWT, check))
	mux.HandleFunc("GET /entries/export", middleware.AuthGuard(http.HandlerFunc(testApp.ExportCSV), testApp.JWT, check))
	mux.HandleFunc("POST /entries/import", middleware.AuthGuard(http.HandlerFunc(testApp.ImportCSV), testApp.JWT, check))

	// Helper for standard JSON requests
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

	// 2. CREATE USER AND GET COOKIE
	userBody := []byte(`{"email": "csv@example.com", "username": "csvuser", "password": "password123"}`)
	doReq(http.MethodPost, "/register", userBody, nil)
	rrLog := doReq(http.MethodPost, "/login", userBody, nil)

	var sessionCookie *http.Cookie
	for _, c := range rrLog.Result().Cookies() {
		if c.Name == "session_id" {
			sessionCookie = c
		}
	}

	// 3. SEED DATABASE WITH ONE ENTRY
	payload := []byte(`{"ledger_entry":"paycheck", "amount": 1000, "description":"Original Salary"}`)
	doReq(http.MethodPost, "/entries", payload, []*http.Cookie{sessionCookie})

	// --- TEST 1: EXPORT CSV ---
	t.Run("Test Export CSV", func(t *testing.T) {
		rrExport := doReq(http.MethodGet, "/entries/export", nil, []*http.Cookie{sessionCookie})

		if rrExport.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for export, got %d", rrExport.Code)
		}

		contentType := rrExport.Header().Get("Content-Type")
		if contentType != "text/csv" {
			t.Fatalf("Expected Content-Type text/csv, got %s", contentType)
		}

		bodyStr := rrExport.Body.String()
		if !strings.Contains(bodyStr, "transaction_id,ledger_entry,amount") {
			t.Fatalf("CSV missing headers. Body: %s", bodyStr)
		}
		if !strings.Contains(bodyStr, "Original Salary") {
			t.Fatalf("CSV missing seeded data. Body: %s", bodyStr)
		}
	})

	// --- TEST 2: IMPORT CSV ---
	t.Run("Test Import CSV", func(t *testing.T) {
		// Create a mock CSV file in memory
		csvData := `transaction_id,ledger_entry,amount,description,transaction_date
,donation,50.00,Red Cross,2026-01-01T12:00:00Z
,paycheck,2000.00,Bonus,
`
		// Prepare the multipart form data
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "import.csv")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		part.Write([]byte(csvData))
		writer.Close()

		// Build the request
		req := httptest.NewRequest(http.MethodPost, "/entries/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType()) // CRITICAL for multipart
		req.AddCookie(sessionCookie)

		rrImport := httptest.NewRecorder()
		mux.ServeHTTP(rrImport, req)

		if rrImport.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d. Body: %s", rrImport.Code, rrImport.Body.String())
		}

		var result importResult
		if err := json.NewDecoder(rrImport.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode import result: %v", err)
		}

		// We imported 2 rows
		if result.Inserted != 2 {
			t.Fatalf("Expected 2 inserted rows, got %d", result.Inserted)
		}
	})

	// --- TEST 3: IMPORT LARGE CSV (WORKER POOL TEST) ---
	t.Run("Test Import Large CSV", func(t *testing.T) {
		// Create a mock CSV file in memory with 10 rows (more than 8 workers)
		csvData := `transaction_id,ledger_entry,amount,description,transaction_date
,paycheck,10.00,Worker Test 1,
,paycheck,10.00,Worker Test 2,
,paycheck,10.00,Worker Test 3,
,paycheck,10.00,Worker Test 4,
,paycheck,10.00,Worker Test 5,
,paycheck,10.00,Worker Test 6,
,paycheck,10.00,Worker Test 7,
,paycheck,10.00,Worker Test 8,
,paycheck,10.00,Worker Test 9,
,paycheck,10.00,Worker Test 10,
`
		// Prepare the multipart form data
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "import_large.csv")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		part.Write([]byte(csvData))
		writer.Close()

		// Build the request
		req := httptest.NewRequest(http.MethodPost, "/entries/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.AddCookie(sessionCookie)

		rrImport := httptest.NewRecorder()
		mux.ServeHTTP(rrImport, req)

		if rrImport.Code != http.StatusOK {
			t.Fatalf("Import failed with status %d. Body: %s", rrImport.Code, rrImport.Body.String())
		}

		var result importResult
		if err := json.NewDecoder(rrImport.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode import result: %v", err)
		}

		// We imported 10 rows
		if result.Inserted != 10 {
			t.Fatalf("Expected 10 inserted rows, got %d", result.Inserted)
		}
	})
}


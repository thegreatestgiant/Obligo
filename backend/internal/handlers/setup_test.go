package handlers

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Global variable to hold our app state during tests
var testApp *App

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			godotenv.Load(envPath)
			log.Printf("Successfully loaded .env from: %s", envPath)
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir { // Reached the root of the hard drive
			break
		}
		dir = parent
	}

	// 1. Connect to the Test Database
	testURL := os.Getenv("DB_TEST_URL")
	if testURL == "" {
		// UPDATE THIS FALLBACK TO MATCH YOUR ACTUAL CREDENTIALS JUST IN CASE!
		// e.g., "postgres://youruser:yourpass@localhost:5435/charity_test?sslmode=disable"
		testURL = "postgres://postgres:postgres@localhost:5435/charity_test?sslmode=disable"
	}

	testDB, err := sql.Open("postgres", testURL)
	if err != nil {
		log.Fatalf("Could not connect to test DB: %v", err)
	}

	// 2. Build the Schema
	// We consolidate your 000, 001, and 002 SQL files here so the test DB is perfectly mapped.
	schemaSQL := `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";
		
		CREATE TABLE IF NOT EXISTS Users (
			user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(255) NOT NULL, 
			password_hash VARCHAR(255) NOT NULL,
			donation_percentage INT DEFAULT 10 
		);
		
		CREATE TABLE IF NOT EXISTS Ledgers (
			transaction_id SERIAL PRIMARY KEY,
			user_id UUID REFERENCES Users(user_id) ON DELETE CASCADE,
			ledger_entry VARCHAR(50) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			description TEXT,
			charity_owed DECIMAL(10,2),
			transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS denylist (
			jti uuid PRIMARY KEY,
			expires_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS refresh_tokens (
			token text PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
			expires_at TIMESTAMP NOT NULL,
			revoked_at TIMESTAMP
		);
	`
	_, err = testDB.Exec(schemaSQL)
	if err != nil {
		log.Fatalf("Could not create test schema: %v", err)
	}

	testApp = &App{
		DB:       testDB,
		JWT:      []byte("test-secret-key-123"),
		Lifetime: time.Hour * 24,
	}

	// 4. Run all the Test... functions in this package
	exitCode := m.Run()

	// 5. Cleanup
	testDB.Close()
	os.Exit(exitCode)
}

// clearDatabase is a helper you will call at the top of every test.
// TRUNCATE empties the tables.
// RESTART IDENTITY resets the SERIAL transaction_id counter back to 1.
// CASCADE ensures related rows (like ledgers for a user) are also wiped.
func clearDatabase() {
	_, err := testApp.DB.Exec(`
		TRUNCATE TABLE 
			Users, 
			Ledgers, 
			refresh_tokens, 
			denylist 
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		panic("Failed to clear database: " + err.Error())
	}
}

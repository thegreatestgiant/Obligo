package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLedgerMath(t *testing.T) {
	clearDatabase()

	insertLedger := "INSERT INTO Ledgers (transaction_id, user_id, ledger_entry, amount, charity_owed, transaction_date) VALUES ($1, $2, $3, $4, $5, $6)"

	// ==========================================
	// 1. STANDARD MATH USER
	// ==========================================
	userID := uuid.New()
	testApp.executeTemplate(context.Background(), "INSERT INTO Users (user_id, email, username, password_hash) VALUES ($1, $2, $3, $4)", userID, "math@example.com", "mathuser", "dummyhash")

	// Insert Paychecks & Donations
	testApp.executeTemplate(context.Background(), insertLedger, "1", userID, "paycheck", 1000.00, 100.00, time.Now())
	testApp.executeTemplate(context.Background(), insertLedger, "2", userID, "paycheck", 500.00, 50.00, time.Now())
	// Fixed: amount is 75, charity_owed is 0.
	testApp.executeTemplate(context.Background(), insertLedger, "3", userID, "donation", 75.00, 0.00, time.Now())

	// Assertions (Flat, no t.Run so it aborts cleanly if it fails)
	if earned := testApp.getAmountEarned(context.Background(), userID); earned != 1500.00 {
		t.Fatalf("Math Error: Expected 1500 earned, got %.2f", earned)
	}
	if owed := testApp.getAmountOwed(context.Background(), userID); owed != 150.00 {
		t.Fatalf("Math Error: Expected 150 owed, got %.2f", owed)
	}
	if donated := testApp.getAmountDonated(context.Background(), userID); donated != 75.00 {
		t.Fatalf("Math Error: Expected 75 donated, got %.2f", donated)
	}
	if fulfilled := testApp.getAmountFulfilled(context.Background(), userID); fulfilled != 50.00 {
		t.Fatalf("Math Error: Expected 50%% fulfilled, got %.2f%%", fulfilled)
	}

	// ==========================================
	// 2. EDGE CASE: GHOST USER (Divide by Zero check)
	// ==========================================
	ghostUser := uuid.New()
	if testApp.getAmountEarned(context.Background(), ghostUser) != 0 || testApp.getAmountOwed(context.Background(), ghostUser) != 0 || testApp.getAmountDonated(context.Background(), ghostUser) != 0 {
		t.Fatalf("Ghost user amounts should all be 0")
	}
	if testApp.getAmountFulfilled(context.Background(), ghostUser) != 0 {
		t.Fatalf("Ghost user fulfilled should safely return 0 without panicking")
	}

	// ==========================================
	// 3. EDGE CASE: NEGATIVE CORRECTIONS
	// ==========================================
	// We create a separate user so we don't mutate the data of the first user!
	correctionUser := uuid.New()
	testApp.executeTemplate(context.Background(), "INSERT INTO Users (user_id, email, username, password_hash) VALUES ($1, $2, $3, $4)", correctionUser, "corr@example.com", "corr", "dummyhash")

	testApp.executeTemplate(context.Background(), insertLedger, "4", correctionUser, "donation", 100.00, 0.00, time.Now())
	testApp.executeTemplate(context.Background(), insertLedger, "5", correctionUser, "donation", -25.00, 0.00, time.Now())

	if donated := testApp.getAmountDonated(context.Background(), correctionUser); donated != 75.00 {
		t.Fatalf("Math Error: Expected 75 donated after correction, got %.2f", donated)
	}
}

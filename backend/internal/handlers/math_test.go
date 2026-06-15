package handlers

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLedgerMath(t *testing.T) {
	// 1. Wipe the test database clean
	clearDatabase()

	// 2. Create a dummy user directly in the database
	userID := uuid.New()
	userQuery := "INSERT INTO Users (user_id, email, username, password_hash) VALUES ($1, $2, $3, $4)"
	testApp.executeTemplate(userQuery, userID, "math@example.com", "mathuser", "dummyhash")

	// 3. Insert known ledger entries directly into the DB to perfectly control the math
	insertLedger := "INSERT INTO Ledgers (transaction_id, user_id, ledger_entry, amount, charity_owed, transaction_date) VALUES ($1, $2, $3, $4, $5, $6)"

	// Paycheck 1: $1000 earned, $100 owed
	testApp.executeTemplate(insertLedger, 1, userID, "paycheck", 1000.00, 100.00, time.Now())

	// Paycheck 2: $500 earned, $50 owed
	testApp.executeTemplate(insertLedger, 2, userID, "paycheck", 500.00, 50.00, time.Now())

	// Donation 1: $75 donated
	testApp.executeTemplate(insertLedger, 3, userID, "donation", 75.00, 0.00, time.Now())

	// ==========================================
	// 4. RUN THE MATH ASSERTIONS
	// ==========================================

	t.Run("Calculates Total Earned", func(t *testing.T) {
		earned := testApp.getAmountEarned(userID)
		if earned != 1500.00 {
			t.Errorf("Math Error: Expected 1500 earned, got %.2f", earned)
		}
	})

	t.Run("Calculates Total Owed", func(t *testing.T) {
		owed := testApp.getAmountOwed(userID)
		if owed != 150.00 {
			t.Errorf("Math Error: Expected 150 owed, got %.2f", owed)
		}
	})

	t.Run("Calculates Total Donated", func(t *testing.T) {
		donated := testApp.getAmountDonated(userID)
		if donated != 75.00 {
			t.Errorf("Math Error: Expected 75 donated, got %.2f", donated)
		}
	})

	t.Run("Calculates Percent Fulfilled", func(t *testing.T) {
		fulfilled := testApp.getAmountFulfilled(userID)
		// Expected: (75 donated / 150 owed) * 100 = 50.0%
		if fulfilled != 50.00 {
			t.Errorf("Math Error: Expected 50%% fulfilled, got %.2f%%", fulfilled)
		}
	})

	// ==========================================
	// 5. TEST EDGE CASES
	// ==========================================
	t.Run("Handles Zero Entries Safely", func(t *testing.T) {
		// Create a brand new UUID that has NO database entries
		ghostUser := uuid.New()

		if testApp.getAmountEarned(ghostUser) != 0 {
			t.Errorf("Ghost user earned should be 0")
		}

		if testApp.getAmountOwed(ghostUser) != 0 {
			t.Errorf("Ghost user owed should be 0")
		}

		if testApp.getAmountDonated(ghostUser) != 0 {
			t.Errorf("Ghost user donated should be 0")
		}

		// This is the most critical edge case:
		// If Owed is 0, does the code try to divide by zero and crash?
		if testApp.getAmountFulfilled(ghostUser) != 0 {
			t.Errorf("Ghost user fulfilled should safely return 0 without panicking")
		}
	})

	// ==========================================
	// 6. TEST NEGATIVE CORRECTIONS
	// ==========================================

	t.Run("Handles Negative Corrections", func(t *testing.T) {
		// Insert a negative donation to simulate a refunded charity payment
		// Using transaction_id 4
		testApp.executeTemplate(insertLedger, 4, userID, "donation", -25.00, 0.00, time.Now())

		donated := testApp.getAmountDonated(userID)
		// They previously donated 75. Minus 25, it should now be 50.
		if donated != 50.00 {
			t.Errorf("Math Error: Expected 50 donated after correction, got %.2f", donated)
		}
	})
}

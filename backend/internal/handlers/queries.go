package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// CREATE (INSERT)
// ============================================================================

func (cfg *App) denyList(jti uuid.UUID) {
	query := "INSERT INTO denylist VALUES ($1, $2)"
	cfg.insertTemplate(query, jti, time.Now().Local().Add(cfg.Lifetime))
}

func (cfg *App) addRefresh(token string, user_id uuid.UUID, expires time.Time) {
	query := "INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)"
	cfg.insertTemplate(query, token, user_id, expires)
}

func (cfg *App) setUser(email, username, passwordHash string) error {
	sqlInsert := "INSERT INTO Users (email,username,password_hash,user_id) VALUES ($1,$2,$3,$4)"
	return cfg.insertTemplate(sqlInsert, email, username, passwordHash, uuid.New())
}

func (cfg *App) insertTemplate(query string, fields ...any) error {
	_, err := cfg.DB.Exec(query, fields)
	if err != nil {
		log.Printf("DB Exec Error for query [%s]: %v", query, err)
		return err
	}
	return nil
}

// ============================================================================
// READ (SELECT)
// ============================================================================

func (cfg *App) queryRowTemplate(query string, dest []any, args ...any) error {
	err := cfg.DB.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Catastrophic DB error in query [%s]: %v", query, err)
		}
		return err
	}
	return nil
}

func (cfg *App) blacklisted(jti uuid.UUID) bool {
	query := "SELECT 1 FROM denylist WHERE jti=$1"
	var placeholder int
	err := cfg.queryRowTemplate(query, []any{&placeholder}, jti)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Not in denylist: %v", jti)
			return false
		}
		return true
	}
	return true
}

func (cfg *App) getRefresh(user_id uuid.UUID) string {
	query := "SELECT token FROM refresh_tokens WHERE user_id=$1 AND expires_at>$2 AND revoked_at IS NULL ORDER BY created_at DESC LIMIT 1"
	var token string
	err := cfg.queryRowTemplate(query, []any{&token}, user_id, time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No such token for user_id: %v", user_id)
		}
		return ""
	}
	log.Printf("Here is the token: %v ", token)
	return token
}

func (cfg *App) getDonationPercent(user_id uuid.UUID) float64 {
	query := "SELECT donation_percentage FROM users WHERE user_id=$1"
	percent := 10.0

	err := cfg.queryRowTemplate(query, []any{&percent}, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return percent
		}
		log.Default().Printf("Couldn't get donation percent, using default: %v", err)
		return 10.0
	}
	log.Printf("Donation Percent: %.2f", percent)
	return percent
}

func (cfg *App) getAmountOwed(user_id uuid.UUID) float64 {
	query := "SELECT SUM(charity_owed) FROM Ledgers WHERE user_id=$1"
	var owed sql.NullFloat64

	if err := cfg.queryRowTemplate(query, []any{&owed}, user_id); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return owed.Float64
		}
		log.Default().Printf("Couldn't sum owed: %v", err)
		return 0
	}

	log.Printf("Total Amount Owed: %.2f", owed.Float64)
	return owed.Float64
}

func (cfg *App) getAmountEarned(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='paycheck'"
	var earned sql.NullFloat64

	err := cfg.queryRowTemplate(query, []any{&earned}, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("You didn't earn anything: %v", user_id)
		}
		return 0.0
	}

	log.Printf("Total Amount Earned: %.2f", earned.Float64)
	return earned.Float64
}

func (cfg *App) getAmountDonated(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='donation'"
	var donated sql.NullFloat64

	err := cfg.queryRowTemplate(query, []any{&donated}, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("You didn't earn anything: %v", user_id)
		}
		return 0.0
	}

	log.Printf("Total Amount Donated: %.2f", donated.Float64)
	return donated.Float64
}

func (cfg *App) getUser(username string) (uuid.UUID, string) {
	sqlQuery := "SELECT user_id,password_hash FROM users WHERE username=$1"
	var user_id uuid.UUID
	var pass string

	err := cfg.queryRowTemplate(sqlQuery, []any{&user_id, &pass}, username)
	if err != nil {
		log.Printf("No such user: %s", username)
		return uuid.Nil, ""
	}
	return user_id, pass
}

func (cfg *App) userExists(email, username string) bool {
	query := "SELECT 1 FROM users WHERE email=$1 OR username=$2"
	// Will return nil if empty, and it doesn't exist

	err := cfg.queryRowTemplate(query, []any{}, email, username)
	cfg.DB.QueryRow(query, email, username).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return true // Fail safe: prevent registration if DB is broken
	}
	return true
}

// Notice we added user_id to the function arguments!
func (cfg *App) getEntry(transaction_id string, user_id uuid.UUID) (Ledger, error) {
	query := `SELECT transaction_id, ledger_entry, amount, COALESCE(description,'') AS description,
	charity_owed, charity_fulfilled 
	FROM Ledgers 
	WHERE transaction_id=$1 AND user_id=$2 
	ORDER BY transaction_date DESC 
	Limit 1`

	var entry Ledger

	// 1. Pack all the pointers into a slice
	dest := []any{
		&entry.TransactionID,
		&entry.LedgerEntry,
		&entry.Amount,
		&entry.Description,
		&entry.CharityOwed,
		&entry.CharityFulfilled,
	}

	// 2. Pass the query, destination slice, and BOTH arguments to the template
	err := cfg.queryRowTemplate(query, dest, transaction_id, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entry, fmt.Errorf("no entry found with ID: %s", transaction_id)
		}
		// The template logged the error, but we still return it so the HTTP handler can send a 500
		return entry, fmt.Errorf("database query failed: %w", err)
	}

	log.Printf("Entry: %v", entry)
	return entry, nil
}

// ============================================================================
// UPDATE
// ============================================================================

func (cfg *App) revokeRefresh(token string) {
	query := "UPDATE refresh_tokens SET updated_at=$1, revoked_at=$1 WHERE token=$2"

	_, err := cfg.DB.Exec(query, time.Now(), token)
	if err != nil {
		log.Printf("Couldn't revoke refresh token: %v", err)
	}
}

// ============================================================================
// DELETE
// ============================================================================

func (cfg *App) deleteTemplate(query string) {
	_, err := cfg.DB.Exec(query)
	if err != nil {
		log.Printf("Couldn't delete: %v", err)
	}
}

func (cfg *App) deleteExpiredJTI() {
	query := "DELETE FROM denylist WHERE expires_at < Now()"
	cfg.deleteTemplate(query)
}

func (cfg *App) deleteExpiredRefresh() {
	query := "DELETE FROM refresh_tokens WHERE expires_at < Now()"
	cfg.deleteTemplate(query)
}

// ============================================================================
// LOGIC WRAPPERS (No direct SQL)
// ============================================================================

func (cfg *App) getAmountFulfilled(user_id uuid.UUID) float64 {
	owed := cfg.getAmountOwed(user_id)
	if owed == 0 {
		return 0.0
	}
	fulfilled := (cfg.getAmountDonated(user_id) / owed) * 100

	log.Printf("Total Percent Fulfilled: %.2f%%", fulfilled)
	return fulfilled
}

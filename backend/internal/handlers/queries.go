package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (cfg *App) executeTemplate(query string, args ...any) (int64, error) {
	result, err := cfg.DB.Exec(query, args...)
	if err != nil {
		log.Printf("DB Exec Error in query [%s]: %v", query, err)
		return 0, err
	}
	return result.RowsAffected()
}

func (cfg *App) queryTemplate(query string, args ...any) (*sql.Rows, error) {
	rows, err := cfg.DB.Query(query, args...)
	if err != nil {
		log.Printf("Catastrophic DB error in query [%s]: %v", query, err)
		return nil, err
	}
	return rows, nil
}

func (cfg *App) queryReturnTemplate(query string, dest []any, args ...any) error {
	err := cfg.DB.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Catastrophic DB error in query [%s]: %v", query, err)
		}
		return err
	}
	return nil
}

func (cfg *App) queryExecTemplate(query string, fields ...any) error {
	_, err := cfg.DB.Exec(query, fields...)
	if err != nil {
		log.Printf("DB Exec Error for query [%s]: %v", query, err)
		return err
	}
	return nil
}

// ============================================================================
// CREATE (INSERT)
// ============================================================================

func (cfg *App) denyList(jti uuid.UUID) {
	query := "INSERT INTO denylist VALUES ($1, $2)"
	cfg.queryExecTemplate(query, jti, time.Now().Local().Add(cfg.Lifetime))
}

func (cfg *App) addRefresh(token string, user_id uuid.UUID, expires time.Time) {
	query := "INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)"
	cfg.queryExecTemplate(query, token, user_id, expires)
}

func (cfg *App) setUser(email, username, passwordHash string) error {
	sqlInsert := "INSERT INTO Users (email,username,password_hash,user_id) VALUES ($1,$2,$3,$4)"
	return cfg.queryExecTemplate(sqlInsert, email, username, passwordHash, uuid.New())
}

func (cfg *App) insertEntry(user_id uuid.UUID, t EntryType, amount float64, description string, owed float64) (newID int, e error) {
	sqlInsert := `INSERT INTO Ledgers 
	(user_id, ledger_entry, amount, description, charity_owed) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING transaction_id`

	err := cfg.queryReturnTemplate(sqlInsert, []any{&newID}, user_id, t, amount, description, owed)
	return newID, err
}

// ============================================================================
// READ (SELECT)
// ============================================================================

func (cfg *App) getDups(user_id uuid.UUID, ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	getRows := `
		SELECT transaction_id 
		FROM Ledgers 
		WHERE user_id = $1 
		AND transaction_id = ANY($2::int[])
	`

	rows, err := cfg.queryTemplate(getRows, user_id, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var duplicateIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		duplicateIDs = append(duplicateIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return duplicateIDs, nil
}

func (cfg *App) getUserEntries(user_id uuid.UUID) ([]Ledger, error) {
	query := `SELECT transaction_id, ledger_entry, amount, COALESCE(description,'') AS description, 
	charity_owed, transaction_date 
	FROM Ledgers 
	WHERE user_id=$1 
	ORDER BY transaction_date DESC`

	rows, err := cfg.queryTemplate(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []Ledger{}

	for rows.Next() {
		var entry Ledger
		err := rows.Scan(
			&entry.TransactionID,
			&entry.LedgerEntry,
			&entry.Amount,
			&entry.Description,
			&entry.CharityOwed,
			&entry.TransactionDate,
		)
		if err != nil {
			log.Printf("Couldn't scan row: %v", err)
			continue
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return entries, nil
}

func (cfg *App) getPass(user_id uuid.UUID) (string, error) {
	getQuery := "SELECT password_hash FROM users WHERE user_id=$1"
	var pass string
	err := cfg.queryReturnTemplate(getQuery, []any{&pass}, user_id)
	if err != nil {
		log.Printf("Bad Password, DB errored: %v", err)
		return pass, err
	}

	return pass, nil
}

func (cfg *App) blacklisted(jti uuid.UUID) bool {
	query := "SELECT 1 FROM denylist WHERE jti=$1"
	var placeholder int
	err := cfg.queryReturnTemplate(query, []any{&placeholder}, jti)
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
	err := cfg.queryReturnTemplate(query, []any{&token}, user_id, time.Now())
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

	err := cfg.queryReturnTemplate(query, []any{&percent}, user_id)
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

	if err := cfg.queryReturnTemplate(query, []any{&owed}, user_id); err != nil {
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

	err := cfg.queryReturnTemplate(query, []any{&earned}, user_id)
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

	err := cfg.queryReturnTemplate(query, []any{&donated}, user_id)
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

	err := cfg.queryReturnTemplate(sqlQuery, []any{&user_id, &pass}, username)
	if err != nil {
		log.Printf("No such user: %s", username)
		return uuid.Nil, ""
	}
	return user_id, pass
}

func (cfg *App) userExists(email, username string) bool {
	query := "SELECT 1 FROM users WHERE email=$1 OR username=$2"
	// Will return nil if empty, and it doesn't exist

	err := cfg.queryReturnTemplate(query, []any{}, email, username)
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
	charity_owed
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
	}

	// 2. Pass the query, destination slice, and BOTH arguments to the template
	err := cfg.queryReturnTemplate(query, dest, transaction_id, user_id)
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

	err := cfg.queryExecTemplate(query, time.Now(), token)
	if err != nil {
		log.Printf("Couldn't revoke refresh token: %v", err)
	}
}

func (cfg *App) updatePercentQuery(percent float64, user_id uuid.UUID) error {
	query := "UPDATE Users SET donation_percentage=$1 WHERE user_id=$2"
	return cfg.queryExecTemplate(query, percent, user_id)
}

func (cfg *App) updatePassword(pass []byte, user_id uuid.UUID) error {
	updateQuery := "UPDATE Users SET password_hash=$1 WHERE user_id=$2"
	return cfg.queryExecTemplate(updateQuery, pass, user_id)
}

func (cfg *App) updateEntry(amount, owed float64, description, transaction_id string, user_id uuid.UUID) error {
	updateQuery := `UPDATE Ledgers SET 
	amount=$1, charity_owed=$2, description=$3 
	WHERE user_id=$4 AND transaction_id=$5`
	return cfg.queryExecTemplate(updateQuery, amount, owed, description, user_id, transaction_id)
}

// ============================================================================
// DELETE
// ============================================================================

func (cfg *App) deleteExpiredJTI() {
	query := "DELETE FROM denylist WHERE expires_at < Now()"
	cfg.queryExecTemplate(query)
}

func (cfg *App) deleteExpiredRefresh() {
	query := "DELETE FROM refresh_tokens WHERE expires_at < Now()"
	cfg.queryExecTemplate(query)
}

func (cfg *App) deleteUserEntry(transaction_id string, user_id uuid.UUID) error {
	query := "DELETE FROM Ledgers WHERE transaction_id=$1 AND user_id=$2"

	rowsAffected, err := cfg.executeTemplate(query, transaction_id, user_id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("not found or unauthorized")
	}

	return nil
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

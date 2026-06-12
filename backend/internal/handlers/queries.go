package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

func (cfg *App) denyList(jti uuid.UUID) {
	query := "INSERT INTO denylist VALUES ($1, $2)"
	_, err := cfg.DB.Exec(query, jti, time.Now().Local().Add(cfg.Lifetime))
	if err != nil {
		log.Printf("Bad jti or time: %v", jti)
	}
}

func (cfg *App) blacklisted(jti uuid.UUID) bool {
	query := "SELECT jti FROM denylist WHERE jti=$1"
	row := cfg.DB.QueryRow(query, jti)
	if err := row.Scan(&uuid.UUID{}); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Not in denylist: %v", jti)
			return false
		}
		log.Default().Printf("Something went wrong: %v", err)
		return false
	}
	return true
}

func (cfg *App) getRefresh(user_id uuid.UUID) string {
	query := "SELECT token FROM refresh_tokens WHERE user_id=$1 AND expires_at>$2 AND revoked_at IS NULL ORDER BY created_at DESC LIMIT 1"
	var token string
	row := cfg.DB.QueryRow(query, user_id, time.Now())
	if err := row.Scan(&token); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No such token for user_id: %v", user_id)
			return ""
		}
		log.Default().Printf("Something went wrong: %v", err)
		return ""
	}
	log.Printf("Here is the token: %v ", token)
	return token
}

func (cfg *App) revokeRefresh(token string) {
	query := "UPDATE refresh_tokens SET updated_at=$1, revoked_at=$1 WHERE token=$2"

	_, err := cfg.DB.Exec(query, time.Now(), token)
	if err != nil {
		log.Printf("Couldn't revoke refresh token: %v", err)
	}
}

func (cfg *App) addRefresh(token string, user_id uuid.UUID, expires time.Time) {
	query := "INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)"
	_, err := cfg.DB.Exec(query, token, user_id, expires)
	if err != nil {
		log.Println("Couldn't creat refresh token")
	}
}

func (cfg *App) getDonationPercent(user_id uuid.UUID) float64 {
	query := "SELECT donation_percentage FROM users WHERE user_id=$1"
	percent := 10.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&percent); err != nil {
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

func (cfg *App) getEntry(user_id uuid.UUID) Ledger {
	query := "SELECT ledger_entry, amount, charity_owed, charity_fulfilled FROM Ledgers WHERE user_id=$1 ORDER BY transaction_date DESC Limit 1"
	var entry Ledger

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&entry.LedgerEntry, &entry.Amount, &entry.CharityOwed, &entry.CharityFulfilled); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v. Or the date was wrong: %v", user_id)
			return entry
		}
		log.Default().Printf("No such Entry: %v", err)
		return entry
	}

	log.Printf("Entry: %v", entry)
	return entry
}

func (cfg *App) getAmountOwed(user_id uuid.UUID) float64 {
	query := "SELECT SUM(charity_owed) FROM Ledgers WHERE user_id=$1"
	var owed sql.NullFloat64

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&owed); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return owed.Float64
		}
		log.Default().Printf("Couldn't sum owed: %v", err)
		return 10
	}

	log.Printf("Total Amount Owed: %.2f", owed)
	return owed.Float64
}

func (cfg *App) getAmountEarned(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='paycheck'"
	var earned sql.NullFloat64
	z := 0.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&earned); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return z
		}
		log.Default().Printf("Couldn't sum fulfilled: %v", err)
		return 10
	}

	log.Printf("Total Amount Earned: %.2f", earned)
	return earned.Float64
}

func (cfg *App) getAmountDonated(user_id uuid.UUID) float64 {
	query := "SELECT SUM(amount) FROM Ledgers WHERE user_id=$1 AND ledger_entry='donation'"
	var donated sql.NullFloat64
	d := 0.0

	row := cfg.DB.QueryRow(query, user_id)
	if err := row.Scan(&donated); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Bad uuid: %v", user_id)
			return d
		}
		log.Default().Printf("Couldn't sum fulfilled: %v", err)
		return 10
	}

	log.Printf("Total Amount Donated: %.2f", donated)
	return donated.Float64
}

func (cfg *App) getAmountFulfilled(user_id uuid.UUID) float64 {
	owed := cfg.getAmountOwed(user_id)
	if owed == 0 {
		return 0.0
	}
	fulfilled := (cfg.getAmountDonated(user_id) / owed) * 100

	log.Printf("Total Percent Fulfilled: %.2f%%", fulfilled)
	return fulfilled
}

func (cfg *App) getUser(username string) (uuid.UUID, string) {
	sqlQuery := "SELECT user_id,password_hash FROM users WHERE username=$1"

	var user_id uuid.UUID
	var pass string
	row := cfg.DB.QueryRow(sqlQuery, username)
	if err := row.Scan(&user_id, &pass); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Incorrect username %s or password", username)
			return uuid.Nil, ""
		}
		log.Printf("Incorrect username %s or password", user_id)
		return uuid.Nil, ""
	}
	return user_id, pass
}

func (cfg *App) setUser(email, username, passwordHash string) error {
	sqlInsert := "INSERT INTO Users (email,username,password_hash,user_id) VALUES ($1,$2,$3,$4)"

	_, err := cfg.DB.Query(sqlInsert, email, username, passwordHash, uuid.New())
	if err != nil {
		return err
	}
	return nil
}

func (cfg *App) userExists(email, username string) bool {
	query := "SELECT * FROM users WHERE email=$1 OR username=$2"
	// Will return nil if empty, and it doesn't exist
	err := cfg.DB.QueryRow(query, email, username).Scan()
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func (cfg *App) deleteExpiredJTI() {
	query := "DELETE FROM denylist WHERE expires_at < Now()"
	_, err := cfg.DB.Exec(query)
	if err != nil {
		log.Printf("Couldn't delete: %v", err)
	}
}

func (cfg *App) deleteExpiredRefresh() {
	query := "DELETE FROM refresh_tokens WHERE expires_at < Now()"
	_, err := cfg.DB.Exec(query)
	if err != nil {
		log.Printf("Couldn't delete: %v", err)
	}
}

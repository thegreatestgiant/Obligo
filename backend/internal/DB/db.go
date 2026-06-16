package DB

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

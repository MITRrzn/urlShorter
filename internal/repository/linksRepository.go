package repository

import (
	"database/sql"
)

func AddLink(db *sql.DB, sourceURL string, shortURL string) error {
	_, err := db.Exec(
		`INSERT INTO links(short_code, original_url)
			VALUES($1, $2)
			`,
		shortURL,
		sourceURL,
	)
	return err
}

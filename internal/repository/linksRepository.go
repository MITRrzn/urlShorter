package repository

import (
	"database/sql"
	"urlShorter/internal/structs"
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

func GetUrlByShortCode(db *sql.DB, shortCode string) (structs.LinkResponse, error) {
	var result structs.LinkResponse

	err := db.QueryRow(
		`
		SELECT
    		l.short_code AS short_code,
    		l.original_url AS original_url

    		FROM links l
    		
			WHERE l.short_code = $1;
		`,
		shortCode,
	).Scan(&result)

	if err != nil {
		return structs.LinkResponse{}, err
	}

	return result, nil
}

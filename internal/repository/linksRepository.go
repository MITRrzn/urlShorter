package repository

import (
	"context"
	"database/sql"
	"urlShorter/internal/structs"
)

func AddLink(ctx context.Context, db *sql.DB, sourceURL string, shortURL string) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO links(short_code, original_url)
			VALUES($1, $2)
			`,
		shortURL,
		sourceURL,
	)
	return err
}

func GetUrlByShortCode(ctx context.Context, db *sql.DB, shortCode string) (structs.LinkResponse, error) {
	var result structs.LinkResponse

	err := db.QueryRowContext(
		ctx,
		`
		SELECT
		    l.id as id,
    		l.short_code AS short_code,
    		l.original_url AS original_url

    		FROM links l
    		
			WHERE l.short_code = $1;
		`,
		shortCode,
	).Scan(&result.ID, &result.ShortURL, &result.OriginalURL)

	if err != nil {
		return structs.LinkResponse{}, err
	}

	return result, nil
}

package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
	"urlShorter/internal/structs"
)

func AddClick(ctx context.Context, db *sql.DB, event structs.ClickEvent) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO clicks(link_id, clicked_at, referer, user_agent, ip, event_id)
    VALUES($1, $2, $3, $4, $5, $6)`,
		event.LinkID,
		event.ClickedAt,
		event.Referer,
		event.UserAgent,
		event.Ip,
		event.EventID,
	)
	return err
}

func GetYesterdayClicks(ctx context.Context, db *sql.DB, previousDayStart time.Time, previousDayEnd time.Time) ([]structs.LinkStat, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT link_id, COUNT(*)
    			FROM clicks
    				WHERE clicked_at >= $1 AND clicked_at < $2
    			GROUP BY link_id
    	`,
		previousDayStart,
		previousDayEnd,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	var results []structs.LinkStat

	for rows.Next() {
		var data structs.LinkStat

		if scanErr := rows.Scan(
			&data.LinkID,
			&data.ClicksCount,
		); scanErr != nil {
			return nil, scanErr
		}

		data.StatDate = previousDayStart

		results = append(results, data)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return results, nil
}

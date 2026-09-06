package repository

import (
	"context"
	"database/sql"
	"time"
	"urlShorter/internal/structs"
)

func GetStatByCode(ctx context.Context, db *sql.DB, code string) (structs.Stats, error) {
	var stat structs.Stats

	err := db.QueryRowContext(
		ctx,
		`
		SELECT
    		l.short_code,
    		COALESCE(SUM(lsd.clicks_count), 0) AS total_clicks
			FROM links l
				LEFT JOIN link_stats_daily lsd ON lsd.link_id = l.id
			WHERE l.short_code = $1
			GROUP BY l.id, l.short_code;`,
		code,
	).Scan(&stat.ShortURL, &stat.ClicksAmount)
	if err != nil {
		return structs.Stats{}, err
	}

	dailyStats, dailyStatsErr := getDailyStats(ctx, db, code)
	if dailyStatsErr != nil {
		return structs.Stats{}, dailyStatsErr
	}
	stat.Daily = dailyStats

	return stat, nil
}

func getDailyStats(
	ctx context.Context,
	db *sql.DB,
	shortCode string,
) ([]structs.DailyStats, error) {
	rows, queryErr := db.QueryContext(
		ctx,
		`
        SELECT
            stat_date,
            clicks_count
        FROM link_stats_daily
        WHERE link_id = (SELECT id FROM links WHERE short_code = $1)
        ORDER BY stat_date
        `,
		shortCode,
	)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var dailyStats []structs.DailyStats
	for rows.Next() {
		var stat structs.DailyStats
		var date time.Time

		if err := rows.Scan(&date, &stat.Clicks); err != nil {
			return nil, err
		}

		stat.Date = date.Format("2006-01-02")
		dailyStats = append(dailyStats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dailyStats, nil
}

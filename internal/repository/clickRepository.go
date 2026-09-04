package repository

import (
	"context"
	"database/sql"
	"urlShorter/internal/structs"
)

func AddClick(ctx context.Context, db *sql.DB, event structs.ClickEvent) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO clicks(link_id, clicked_at, referer, user_agent, ip_hash, event_id)
    VALUES($1, $2, $3, $4, $5)`,
		event.LinkID,
		event.ClickedAt,
		event.Referer,
		event.UserAgent,
		event.Ip,
		event.EventID,
	)
	return err
}

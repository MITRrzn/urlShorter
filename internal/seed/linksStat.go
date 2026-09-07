package seed

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"

	"github.com/schollz/progressbar/v3"
)

func LinksStatSeed(db *sql.DB) int64 {
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `
		SELECT id
		FROM links
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		rowsCloseErr := rows.Close()
		if rowsCloseErr != nil {
			return
		}
	}(rows)

	var linkIDs []int64

	for rows.Next() {
		var id int64
		if scanErr := rows.Scan(&id); scanErr != nil {
			log.Fatal(scanErr)
		}

		linkIDs = append(linkIDs, id)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		log.Fatal(rowsErr)
	}

	now := time.Now().UTC()

	var total int64
	pb := progressbar.Default(int64(len(linkIDs)), "linksStat seeding")
	for _, linkID := range linkIDs {
		pbErr := pb.Add(1)
		if pbErr != nil {
			continue
		}
		daysCount := rand.Intn(21) + 10

		for i := 0; i < daysCount; i++ {
			statDate := now.AddDate(0, 0, -i)
			clicksCount := rand.Intn(200) + 1

			addErr := repository.AddDailyLinkStat(ctx, db, structs.LinkStat{
				LinkID:      linkID,
				StatDate:    statDate,
				ClicksCount: int64(clicksCount),
			})
			if addErr != nil {
				log.Fatal(addErr)
			}
			total++
		}
	}
	return total
}

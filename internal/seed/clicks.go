package seed

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
)

func SeedClicks(db *sql.DB) int64 {
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `
		SELECT id
		FROM links
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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

	var total int64
	pb := progressbar.Default(int64(len(linkIDs)))
	for _, linkID := range linkIDs {
		pbErr := pb.Add(1)
		if pbErr != nil {
			continue
		}
		clicksCount := rand.Intn(141) + 10

		for i := 0; i < clicksCount; i++ {

			_, dbExecErr := db.ExecContext(
				ctx,
				`
				INSERT INTO clicks (
					link_id,
					clicked_at,
					referer,
					user_agent,
					ip_hash,
					event_id
				)
				VALUES ($1, $2, $3, $4, $5, $6)
				`,
				linkID,
				randomDateLast30Days(),
				faker.URL(),
				fakeUserAgent(),
				faker.IPv4(),
				uuid.New(),
			)
			if dbExecErr != nil {
				log.Fatal(dbExecErr)
			}
			total++
		}
	}

	return total
}

func randomDateLast30Days() time.Time {
	now := time.Now().UTC()

	randomSeconds := rand.Int63n(int64(30 * 24 * time.Hour / time.Second))

	return now.Add(-time.Duration(randomSeconds) * time.Second)
}

func fakeUserAgent() string {
	generator := faker.GetUserAgent()

	value, err := generator.UserAgent(reflect.Value{})
	if err != nil {
		return ""
	}

	userAgent, ok := value.(string)
	if !ok {
		return ""
	}

	return userAgent
}

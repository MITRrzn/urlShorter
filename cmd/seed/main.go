package main

import (
	"database/sql"
	"flag"
	"fmt"
	"urlShorter/internal/database"
	"urlShorter/internal/seed"
)

func main() {
	amount := flag.Int64("links", 100, "amount of links to generate")
	flag.Parse()

	db, err := database.PsqlConnect()
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		closeErr := db.Close()
		if closeErr != nil {

		}
	}(db)

	totalLinks := seed.SeedLinks(db, *amount)
	totalClicks := seed.SeedClicks(db)

	fmt.Println(fmt.Sprintf("%d links seeded", totalLinks))
	fmt.Println(fmt.Sprintf("%d clicks seeded", totalClicks))
}

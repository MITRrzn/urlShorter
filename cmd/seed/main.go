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

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		closeErr := db.Close()
		if closeErr != nil {

		}
	}(db)

	total := seed.Run(db, *amount)
	fmt.Println(fmt.Sprintf("%d links seeded", total))
}

package seed

import (
	"database/sql"
	"urlShorter/internal/repository"
	"urlShorter/internal/shorter"

	"github.com/go-faker/faker/v4"
	"github.com/schollz/progressbar/v3"
)

func Run(db *sql.DB, amount int64) int64 {
	var total int64

	pb := progressbar.Default(amount)
	for i := 0; int64(i) < amount; i++ {
		err := pb.Add(1)
		sourceURL := faker.URL()
		shortURL, err := shorter.GenerateShortUrl()
		if err != nil {
			continue
		}

		err = repository.AddLink(db, sourceURL, shortURL)
		if err != nil {
			continue
		}
		total++
	}

	return total
}

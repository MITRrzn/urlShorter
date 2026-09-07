package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"urlShorter/internal/database"
	"urlShorter/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.PsqlConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			log.Printf("db close error: %v", dbCloseErr)
		}
	}()

	now := time.Now().UTC()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	previousDayStart := todayStart.AddDate(0, 0, -1)
	previousDayEnd := todayStart

	prevDayData, repoErr := repository.GetYesterdayClicks(ctx, db, previousDayStart, previousDayEnd)
	if repoErr != nil {
		log.Fatal(repoErr)
	}

	for _, data := range prevDayData {
		addErr := repository.AddDailyLinkStat(ctx, db, data)
		if addErr != nil {
			log.Printf("failed to add daily stat: %v", addErr)
			return
		}
	}
}

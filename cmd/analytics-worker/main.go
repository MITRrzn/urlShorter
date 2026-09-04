package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"urlShorter/internal/database"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"

	"github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   "urlShorter-clicks",
		GroupID: "analytics-worker",
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Println(err)
		}
	}(reader)

	db, psqlErr := database.PsqlConnect()
	if psqlErr != nil {
		log.Fatal(psqlErr)
	}
	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			log.Printf("db close error: %v", dbCloseErr)
		}
	}()

	for {
		var event structs.ClickEvent

		msg, fetchErr := reader.FetchMessage(ctx)
		if errors.Is(fetchErr, context.Canceled) {
			break
		}
		if fetchErr != nil {
			log.Println(fetchErr)
			continue
		}

		unmarshalErr := json.Unmarshal(msg.Value, &event)
		if unmarshalErr != nil {
			commitErr := reader.CommitMessages(ctx, msg)
			if commitErr != nil {
				log.Println("Error committing bad message", commitErr)
				return
			}

			log.Println("Error unmarshalling", unmarshalErr)
			continue
		}

		addErr := repository.AddClick(ctx, db, event)
		if addErr != nil {
			var pqErr *pq.Error

			isDuplicateEvent := errors.As(addErr, &pqErr) && pqErr.Code == "23505" && pqErr.Constraint == "clicks_event_id_key"
			if !isDuplicateEvent {
				log.Println("Error adding click", addErr)
				continue
			}

			log.Println("duplicate event:", event.EventID)
		}

		commitErr := reader.CommitMessages(ctx, msg)
		if commitErr != nil {
			log.Println("Error committing message", commitErr)
			return
		}
	}
}

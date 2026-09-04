package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlShorter/internal/database"
	kafkaClient "urlShorter/internal/kafka"
	"urlShorter/internal/redirect"
	"urlShorter/internal/shorter"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	db, err := database.PsqlConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			log.Printf("db close error: %v", dbCloseErr)
		}
	}()

	redisClient, redisErr := database.GetRedisClient(ctx)
	if redisErr != nil {
		log.Fatal(redisErr)
	}

	writer := kafkaClient.GetWriter()
	defer func(writer *kafka.Writer) {
		writerErr := writer.Close()
		if writerErr != nil {
			log.Println(writerErr)
		}
	}(writer)

	mux.HandleFunc("POST /links", shorter.CreateLinkHandler(db))
	mux.HandleFunc("GET /{code}", redirect.RedirectHandler(db, redisClient, writer))

	port := os.Getenv("APP_PORT")
	log.Println("Starting server at port", port)
	server := &http.Server{
		Addr:              fmt.Sprint(":", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server")
	shutdownErr := server.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		log.Printf("server shutdown error: %v", shutdownErr)
	}
	log.Println("Server stopped")
}

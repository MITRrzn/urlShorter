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
	"urlShorter/internal/redirect"
	"urlShorter/internal/shorter"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			log.Printf("db close error: %v", dbCloseErr)
		}
	}()

	mux.HandleFunc("POST /links", shorter.CreateLinkHandler(db))
	mux.HandleFunc("GET /{code}", redirect.RedirectHandler(db))

	log.Println("Starting server at port 8080")
	server := &http.Server{
		Addr:    fmt.Sprint(":", os.Getenv("APP_PORT")),
		Handler: mux,
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

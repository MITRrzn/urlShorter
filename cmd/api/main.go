package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"urlShorter/internal/redirect"
	"urlShorter/internal/shorter"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", shorter.CreateLinkHandler)
	mux.HandleFunc("GET /{code}", redirect.RedirectHandler)

	log.Println("Starting server at port 8080")
	server := &http.Server{
		Addr:    ":8080",
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

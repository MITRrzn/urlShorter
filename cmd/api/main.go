package main

import (
	"fmt"
	"net/http"
	"time"
	"urlShorter/internal/redirect"
	"urlShorter/internal/shorter"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /links", shorter.CreateLinkHandler)
	mux.HandleFunc("GET /{code}", redirect.RedirectHandler)

	fmt.Printf("%s :: Starting server at port 8080", time.Now().Format("2006-01-02 15:04:05"))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

package shorter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"

	"github.com/lib/pq"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateLinkHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link structs.LinkStruct

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&link); err != nil {
			log.Println("decode error:", err)
			helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
			return
		}

		validateErr := validateURL(link.URL)
		if validateErr != nil {
			helper.WriteErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
			return
		}

		storeResult, storeErr := storeShortURL(r.Context(), db, link)

		if storeErr != nil {
			helper.WriteErrorResponse(w, "failed to create short url", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(structs.SuccessResponse{
			Success:  "created",
			ShortURL: storeResult,
		})
		if encodeErr != nil {
			log.Println(encodeErr)
			return
		}
	}
}

func storeShortURL(ctx context.Context, db *sql.DB, link structs.LinkStruct) (string, error) {
	for i := 0; i < 5; i++ {
		var pqErr *pq.Error

		shortURL, genErr := GenerateShortUrl()
		if genErr != nil {
			log.Println(genErr)
			return "", genErr
		}

		err := repository.AddLink(ctx, db, link.URL, shortURL)
		if err == nil {
			return shortURL, nil
		}
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			continue
		}

		log.Println(err)
		return "", err
	}

	return "", errors.New("failed to generate unique short url")
}

func validateURL(linkURL string) error {
	if linkURL == "" {
		return errors.New("url is empty")
	}

	parsedURL, err := url.Parse(linkURL)
	if err != nil {
		return err
	}

	switch {
	case parsedURL.Scheme == "":
		return errors.New("scheme is empty")
	case parsedURL.Scheme != "http" && parsedURL.Scheme != "https":
		return errors.New("url scheme must be http or https")
	case parsedURL.Host == "":
		return errors.New("host is empty")
	default:
		return nil
	}
}

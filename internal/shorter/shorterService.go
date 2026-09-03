package shorter

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"
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

		shortURL, err := GenerateShortUrl()
		if err != nil {
			helper.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		linkRepoErr := repository.AddLink(db, link.URL, shortURL)
		if linkRepoErr != nil {
			helper.WriteErrorResponse(w, linkRepoErr.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(structs.SuccessResponse{
			Success: "created",
		})
		if encodeErr != nil {
			log.Println(encodeErr)
			return
		}
	}
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

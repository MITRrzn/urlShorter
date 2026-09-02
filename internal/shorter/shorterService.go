package shorter

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"urlShorter/internal/database"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	var link structs.LinkStruct

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&link); err != nil {
		log.Println("decode error:", err)
		writeErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	validateErr := validateURL(link.URL)
	if validateErr != nil {
		writeErrorResponse(w, validateErr.Error(), http.StatusBadRequest)
		return
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

	shortURL, err := generateShortUrl()
	if err != nil {
		writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	linkRepoErr := repository.AddLink(db, link.URL, shortURL)
	if linkRepoErr != nil {
		writeErrorResponse(w, linkRepoErr.Error(), http.StatusInternalServerError)
	}
}

func generateShortUrl() (string, error) {
	result := make([]byte, 7)

	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
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

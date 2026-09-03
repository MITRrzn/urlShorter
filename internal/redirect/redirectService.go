package redirect

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
	"urlShorter/internal/structs"
)

func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		redirectData, err := repository.GetUrlByShortCode(db, code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				helper.WriteErrorResponse(w, "link not found", http.StatusNotFound)
				return
			}

			log.Println(err)
			helper.WriteErrorResponse(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encodeErr := json.NewEncoder(w).Encode(structs.LinkResponse{
			ShortURL:    redirectData.ShortURL,
			OriginalURL: redirectData.OriginalURL,
		})
		if encodeErr != nil {
			log.Println(encodeErr)
			return
		}
	}
}

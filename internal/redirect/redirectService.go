package redirect

import (
	"database/sql"
	"encoding/json"
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
			log.Println(err)
			helper.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encodeErr := json.NewEncoder(w).Encode(structs.LinkResponse{
			ShortURL:    redirectData.OriginalURL,
			OriginalURL: redirectData.OriginalURL,
		})
		if encodeErr != nil {
			log.Println(encodeErr)
			return
		}
	}
}

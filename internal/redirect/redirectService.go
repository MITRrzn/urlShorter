package redirect

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
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

		http.Redirect(w, r, redirectData.OriginalURL, http.StatusFound)
	}
}

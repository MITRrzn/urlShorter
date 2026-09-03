package redirect

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"regexp"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
)

func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		validationErr := validateCode(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		redirectData, err := repository.GetUrlByShortCode(r.Context(), db, code)
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

func validateCode(code string) error {
	if code == "" {
		return errors.New("empty code")
	}
	if len(code) != 7 {
		return errors.New("invalid code")
	}
	var shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{7}$`)
	if !shortCodeRegex.MatchString(code) {
		return errors.New("incorrect code")
	}

	return nil
}

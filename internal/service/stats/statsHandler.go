package stats

import (
	"database/sql"
	"net/http"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
)

func StatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		validationErr := helper.ValidateCode(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		stats, err := repository.GetStatByCode(r.Context(), db, code)
		if err != nil {
			helper.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}

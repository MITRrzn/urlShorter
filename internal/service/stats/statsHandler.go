package stats

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"urlShorter/internal/helper"
	"urlShorter/internal/repository"
)

func StatsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		validationErr := helper.ValidateCode(code)
		if validationErr != nil {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		stats, err := repository.GetStatByCode(r.Context(), db, code)
		if errors.Is(err, sql.ErrNoRows) {
			helper.WriteErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			helper.WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(stats)
		if encodeErr != nil {
			log.Println(encodeErr)
			return
		}
		return
	}
}

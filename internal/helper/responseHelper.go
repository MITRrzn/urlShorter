package helper

import (
	"encoding/json"
	"net/http"
	"urlShorter/internal/structs"
)

func WriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(structs.ErrorResponse{
		Error: errorMessage,
	})
}

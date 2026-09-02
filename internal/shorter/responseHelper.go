package shorter

import (
	"encoding/json"
	"net/http"
	"urlShorter/internal/structs"
)

func writeErrorResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(structs.ErrorResponse{
		Error: errorMessage,
	})
}

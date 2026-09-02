package shorter

import (
	"encoding/json"
	"log"
	"net/http"
	"urlShorter/internal/structs"
)

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link structs.LinkStruct

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&link); err != nil {
		log.Println("decode error:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("created"))
	if err != nil {
		return
	}
}

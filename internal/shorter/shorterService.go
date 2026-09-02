package shorter

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"urlShorter/internal/structs"
	"urlShorter/internal/structs/createLinkResponse"
)

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link structs.LinkStruct

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&link); err != nil {
		log.Println("decode error:", err)
		writeErrorResponse(w, "incorrect json format")
		return
	}

	validateErr := validateURL(link.URL)
	if validateErr != nil {
		writeErrorResponse(w, validateErr.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(createLinkResponse.SuccessResponse{
		Success: "created",
	})
	if encodeErr != nil {
		log.Println(encodeErr)
		return
	}
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

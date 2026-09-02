package shorter

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"urlShorter/internal/structs"
)

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link structs.LinkStruct

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&link); err != nil {
		log.Println("decode error:", err)
		return
	}

	validateErr := validateLink(link)
	if validateErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(validateErr.Error()))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("created"))
	if err != nil {
		return
	}
}

func validateLink(link structs.LinkStruct) error {
	if link.URL == "" {
		return errors.New("url is empty")
	}

	parsedUrl, err := url.Parse(link.URL)
	if err != nil {
		return err
	}

	switch {
	case parsedUrl.Scheme == "":
		return errors.New("scheme is empty")
	case parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https":
		return errors.New("url scheme must be http or https")
	case parsedUrl.Host == "":
		return errors.New("host is empty")
	default:
		return nil
	}
}

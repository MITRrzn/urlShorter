package shorter

import "net/http"

func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("created"))
	if err != nil {
		return
	}
}

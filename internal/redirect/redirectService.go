package redirect

import "net/http"

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	//code := r.PathValue("code")
	tmpRedirectUrlFromDB := "http://google.com"
	http.Redirect(w, r, tmpRedirectUrlFromDB, http.StatusFound)
}

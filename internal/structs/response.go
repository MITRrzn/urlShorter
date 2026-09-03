package structs

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success  string `json:"success"`
	ShortURL string `json:"short_url"`
}

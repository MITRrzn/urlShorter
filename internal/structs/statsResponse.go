package structs

type Stats struct {
	ShortURL     string       `json:"short_url"`
	ClicksAmount int64        `json:"clicks_amount"`
	Daily        []DailyStats `json:"daily"`
}

type DailyStats struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

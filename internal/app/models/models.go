package models

// Request to create short url
type Request struct {
	URL string `json:"url"`
}

// Response with short url
type Response struct {
	Result string `json:"result"`
}

// batch record
type BatchRequestRecord struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// batch records
type BatchRequest []BatchRequestRecord

// batch response record
type BatchResponseRecord struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// batch response
type BatchResponse []BatchResponseRecord

// url record
type URLRecord struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// statictics record
type StatRecord struct {
	URLS  int `json:"urls"`
	Users int `json:"users"`
}

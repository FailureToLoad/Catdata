package domain

type CatRecord struct {
	ID        int     `json:"id"`
	Timestamp int64   `json:"timestamp"`
	Cat       string  `json:"cat"`
	Weight    float32 `json:"weight"`
	Notes     *string `json:"notes"`
}

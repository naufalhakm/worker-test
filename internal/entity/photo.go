package entity

import "time"

type Photo struct {
	ID          uint64  `gorm:"primaryKey" json:"id"`
	UserID      uint64  `json:"user_id"`
	URL         string  `json:"url"`
	Status      string  `json:"status"`
	Similiarity float64 `json:"similiarity"`
	EmbedData   []byte
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

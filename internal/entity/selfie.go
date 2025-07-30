package entity

import "time"

type Selfie struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	UserID    uint64 `json:"user_id"`
	URL       string `json:"url"`
	EmbedData []byte
	CreatedAt time.Time `json:"created_at"`
}

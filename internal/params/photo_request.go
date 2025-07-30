package params

type UploadPhotoRequest struct {
	UserID   uint64 `json:"user_id" validate:"required"`
	PhotoURL string `json:"photo_url" validate:"required"`
}

type UpdatePhotoRequest struct {
	UserID uint64 `json:"user_id" validate:"required"`
	Status string `json:"status"`
}

package params

type UploadSelfieRequest struct {
	UserID    uint64 `json:"user_id" validate:"required"`
	SelfieURL string `json:"selfie_url" validate:"required"`
}

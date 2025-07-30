package handler

import (
	"go-worker/internal/params"
	"go-worker/internal/service"

	"github.com/gin-gonic/gin"
)

type SelfieHandler interface {
	Upload(c *gin.Context)
}
type SelfieHandlerImpl struct {
	SelfieService service.SelfieService
}

func NewSelfieHandler(servicePhoto service.SelfieService) SelfieHandler {
	return &SelfieHandlerImpl{
		SelfieService: servicePhoto,
	}
}

func (h *SelfieHandlerImpl) Upload(c *gin.Context) {
	var req params.UploadSelfieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	err := h.SelfieService.UploadSelfie(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed"})
		return
	}
	c.JSON(200, gin.H{"message": "uploaded"})
}

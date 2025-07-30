package handler

import (
	"go-worker/internal/params"
	"go-worker/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoHandler interface {
	UploadPhoto(c *gin.Context)
	UploadSelfie(c *gin.Context)
}
type PhotoHandlerImpl struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(servicePhoto service.PhotoService) PhotoHandler {
	return &PhotoHandlerImpl{
		PhotoService: servicePhoto,
	}
}

func (h *PhotoHandlerImpl) UploadPhoto(c *gin.Context) {
	var req params.UploadPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	err := h.PhotoService.UploadPhoto(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "uploaded"})
}

func (h *PhotoHandlerImpl) UploadSelfie(c *gin.Context) {
	var req params.UploadSelfieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	err := h.PhotoService.UploadSelfie(c.Request.Context(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "uploaded"})
}

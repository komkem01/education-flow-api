package pictures

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type PresignUploadRequest struct {
	ObjectKey      string `json:"object_key" binding:"required"`
	ExpiresSeconds *int64 `json:"expires_seconds"`
}

type PresignURLResponse struct {
	URL string `json:"url"`
}

func (c *Controller) PresignUpload(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req PresignUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	presignedURL, err := c.svc.PresignUploadURL(ctx.Request.Context(), req.ObjectKey, req.ExpiresSeconds)
	if err != nil {
		c.handleServiceError(ctx, log, err, "picture-presign-upload-failed")
		return
	}

	base.Success(ctx, &PresignURLResponse{URL: presignedURL}, "success")
}

func (c *Controller) PicturesPresignUpload(ctx *gin.Context) {
	c.PresignUpload(ctx)
}

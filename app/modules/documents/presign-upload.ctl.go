package documents

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type PresignUploadRequest struct {
	ObjectKey      string `json:"object_key" binding:"required"`
	Bucket         string `json:"bucket"`
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

	var bucket *string
	if req.Bucket != "" {
		bucket = &req.Bucket
	}

	presignedURL, err := c.svc.PresignUploadURL(ctx.Request.Context(), req.ObjectKey, bucket, req.ExpiresSeconds)
	if err != nil {
		c.handleServiceError(ctx, log, err, "document-presign-upload-failed")
		return
	}

	base.Success(ctx, &PresignURLResponse{URL: presignedURL}, "success")
}

func (c *Controller) DocumentsPresignUpload(ctx *gin.Context) {
	c.PresignUpload(ctx)
}

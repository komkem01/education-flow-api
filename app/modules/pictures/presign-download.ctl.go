package pictures

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type PresignDownloadRequest struct {
	ObjectKey      string `form:"object_key" binding:"required"`
	ExpiresSeconds *int64 `form:"expires_seconds"`
}

func (c *Controller) PresignDownload(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req PresignDownloadRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	presignedURL, err := c.svc.PresignDownloadURL(ctx.Request.Context(), req.ObjectKey, req.ExpiresSeconds)
	if err != nil {
		c.handleServiceError(ctx, log, err, "picture-presign-download-failed")
		return
	}

	base.Success(ctx, &PresignURLResponse{URL: presignedURL}, "success")
}

func (c *Controller) PicturesPresignDownload(ctx *gin.Context) {
	c.PresignDownload(ctx)
}

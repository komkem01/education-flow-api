package documents

import (
	"strings"

	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) Upload(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser == nil || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	documentID, err := uuid.Parse(strings.TrimSpace(ctx.Param("id")))
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	opened, err := fileHeader.Open()
	if err != nil {
		base.InternalServerError(ctx, "document-upload-open-file-failed", nil)
		return
	}
	defer opened.Close()

	contentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	item, err := c.svc.UploadByID(
		ctx.Request.Context(),
		documentID,
		currentUser.Member.SchoolID,
		currentUser.Member.ID,
		currentUser.Member.Role,
		strings.TrimSpace(fileHeader.Filename),
		contentType,
		fileHeader.Size,
		opened,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "document-upload-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) DocumentsUpload(ctx *gin.Context) {
	c.Upload(ctx)
}

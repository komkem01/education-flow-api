package pictures

import (
	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	OwnerMemberID string         `json:"owner_member_id"`
	ObjectKey     string         `json:"object_key" binding:"required"`
	FileName      string         `json:"file_name" binding:"required"`
	ContentType   string         `json:"content_type" binding:"required"`
	SizeBytes     int64          `json:"size_bytes" binding:"required"`
	Metadata      map[string]any `json:"metadata"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser == nil || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request-form", nil)
		return
	}

	var ownerMemberID *uuid.UUID
	if req.OwnerMemberID != "" {
		parsed, err := uuid.Parse(req.OwnerMemberID)
		if err != nil {
			base.BadRequest(ctx, "invalid-request-form", nil)
			return
		}
		if !c.isElevatedRole(currentUser.Member.Role) && parsed != currentUser.Member.ID {
			base.Forbidden(ctx, "unauthorized", nil)
			return
		}
		ownerMemberID = &parsed
	}

	item, err := c.svc.Create(
		ctx.Request.Context(),
		currentUser.Member.SchoolID,
		currentUser.Member.ID,
		ownerMemberID,
		req.ObjectKey,
		req.FileName,
		req.ContentType,
		req.SizeBytes,
		req.Metadata,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "picture-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreatePictureController(ctx *gin.Context) {
	c.Create(ctx)
}

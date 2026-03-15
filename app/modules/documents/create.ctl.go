package documents

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	OwnerMemberID string         `json:"owner_member_id"`
	BucketName    string         `json:"bucket_name"`
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
		ownerMemberID = &parsed
	}

	var bucketName *string
	if req.BucketName != "" {
		if currentUser.Member.Role != ent.MemberRoleAdmin && currentUser.Member.Role != ent.MemberRoleSuperadmin {
			base.Forbidden(ctx, "unauthorized", nil)
			return
		}
		bucketName = &req.BucketName
	}

	item, err := c.svc.Create(
		ctx.Request.Context(),
		currentUser.Member.SchoolID,
		currentUser.Member.ID,
		ownerMemberID,
		bucketName,
		req.ObjectKey,
		req.FileName,
		req.ContentType,
		req.SizeBytes,
		req.Metadata,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "document-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateDocumentController(ctx *gin.Context) {
	c.Create(ctx)
}

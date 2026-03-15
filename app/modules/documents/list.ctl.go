package documents

import (
	"eduflow/app/modules/auth"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ListRequest struct {
	base.RequestPaginate
	OwnerMemberID string `form:"owner_member_id"`
	Status        string `form:"status"`
	StorageID     string `form:"storage_id"`
}

func (c *Controller) List(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser == nil || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
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
	if !c.isElevatedRole(currentUser.Member.Role) {
		if ownerMemberID != nil && *ownerMemberID != currentUser.Member.ID {
			base.Forbidden(ctx, "unauthorized", nil)
			return
		}
		ownerMemberID = &currentUser.Member.ID
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	var storageID *uuid.UUID
	if req.StorageID != "" {
		parsed, err := uuid.Parse(req.StorageID)
		if err != nil {
			base.BadRequest(ctx, "invalid-request-form", nil)
			return
		}
		storageID = &parsed
	}

	items, page, err := c.svc.List(ctx.Request.Context(), &req.RequestPaginate, currentUser.Member.SchoolID, ownerMemberID, status, storageID)
	if err != nil {
		c.handleServiceError(ctx, log, err, "document-list-failed")
		return
	}

	base.Paginate(ctx, items, page)
}

func (c *Controller) DocumentsList(ctx *gin.Context) {
	c.List(ctx)
}

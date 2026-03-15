package members

import (
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	SchoolID  *string `json:"school_id"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
	IsActive  *bool   `json:"is_active"`
	LastLogin *string `json:"last_login"`
}

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var uriReq UpdateByIDRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	id, err := uuid.Parse(uriReq.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var schoolID *uuid.UUID
	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-id", nil)
			return
		}
		schoolID = &parsed
	}

	var role *ent.MemberRole
	if req.Role != nil {
		parsed, ok := parseMemberRole(*req.Role)
		if !ok {
			base.BadRequest(ctx, "invalid-role", nil)
			return
		}
		role = &parsed
	}

	var lastLogin *time.Time
	if req.LastLogin != nil {
		parsed, err := time.Parse(time.RFC3339, *req.LastLogin)
		if err != nil {
			base.BadRequest(ctx, "invalid-last-login", nil)
			return
		}
		lastLogin = &parsed
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, schoolID, req.Email, req.Password, role, req.IsActive, lastLogin)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-update-failed")
		return
	}

	base.Success(ctx, toMemberResponse(item), "success")
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) MembersUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

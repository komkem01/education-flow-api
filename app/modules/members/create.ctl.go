package members

import (
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID  string `json:"school_id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	LastLogin string `json:"last_login"`
}

type memberResponse struct {
	ID        uuid.UUID      `json:"id"`
	SchoolID  uuid.UUID      `json:"school_id"`
	Email     string         `json:"email"`
	Role      ent.MemberRole `json:"role"`
	IsActive  bool           `json:"is_active"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func toMemberResponse(item *ent.Member) *memberResponse {
	if item == nil {
		return nil
	}
	return &memberResponse{
		ID:        item.ID,
		SchoolID:  item.SchoolID,
		Email:     item.Email,
		Role:      item.Role,
		IsActive:  item.IsActive,
		LastLogin: item.LastLogin,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func parseMemberRole(role string) (ent.MemberRole, bool) {
	switch ent.MemberRole(role) {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher, ent.MemberRoleStudent:
		return ent.MemberRole(role), true
	default:
		return "", false
	}
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return
	}

	var role ent.MemberRole
	if req.Role != "" {
		parsed, ok := parseMemberRole(req.Role)
		if !ok {
			base.BadRequest(ctx, "invalid-role", nil)
			return
		}
		role = parsed
	}

	var lastLogin *time.Time
	if req.LastLogin != "" {
		parsed, err := time.Parse(time.RFC3339, req.LastLogin)
		if err != nil {
			base.BadRequest(ctx, "invalid-last-login", nil)
			return
		}
		lastLogin = &parsed
	}

	item, err := c.svc.Create(ctx.Request.Context(), schoolID, req.Email, req.Password, role, req.IsActive, lastLogin)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-create-failed")
		return
	}

	base.Success(ctx, toMemberResponse(item), "success")
}

func (c *Controller) CreateMemberController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) MembersCreate(ctx *gin.Context) {
	c.Create(ctx)
}

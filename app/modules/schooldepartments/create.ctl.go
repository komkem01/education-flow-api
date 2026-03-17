package schooldepartments

import (
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID     string  `json:"school_id"`
	DepartmentID string  `json:"department_id" binding:"required"`
	Code         string  `json:"code"`
	CustomName   *string `json:"custom_name"`
	IsActive     *bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}
	if !canWriteSchoolDepartments(user.Member.Role) {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID := user.Member.SchoolID
	if req.SchoolID != "" {
		parsed, err := uuid.Parse(req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = parsed
	}
	if user.Member.Role != ent.MemberRoleSuperadmin && schoolID != user.Member.SchoolID {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-department-id", nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &ent.SchoolDepartment{
		SchoolID:     schoolID,
		DepartmentID: departmentID,
		Code:         req.Code,
		CustomName:   req.CustomName,
		IsActive:     isActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-department-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateSchoolDepartmentController(ctx *gin.Context) {
	c.Create(ctx)
}

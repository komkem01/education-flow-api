package membermanagements

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	MemberID      string `json:"member_id" binding:"required"`
	EmployeeCode  string `json:"employee_code" binding:"required"`
	Position      string `json:"position" binding:"required"`
	StartWorkDate string `json:"start_work_date" binding:"required"`
	DepartmentID  string `json:"department_id" binding:"required"`
	IsActive      bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	memberID, err := uuid.Parse(req.MemberID)
	if err != nil {
		base.BadRequest(ctx, "invalid-member-id", nil)
		return
	}
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-department-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), memberID, req.EmployeeCode, req.Position, req.StartWorkDate, departmentID, req.IsActive)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-management-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateMemberManagementController(ctx *gin.Context) {
	c.Create(ctx)
}

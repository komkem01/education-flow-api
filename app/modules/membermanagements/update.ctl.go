package membermanagements

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"strings"

	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	MemberID                *string `json:"member_id"`
	MemberIDCamel           *string `json:"memberId"`
	EmployeeCode            *string `json:"employee_code"`
	EmployeeCodeCamel       *string `json:"employeeCode"`
	GenderID                *string `json:"gender_id"`
	GenderIDCamel           *string `json:"genderId"`
	PrefixID                *string `json:"prefix_id"`
	PrefixIDCamel           *string `json:"prefixId"`
	FirstName               *string `json:"first_name"`
	FirstNameCamel          *string `json:"firstName"`
	LastName                *string `json:"last_name"`
	LastNameCamel           *string `json:"lastName"`
	Phone                   *string `json:"phone"`
	PhoneCamel              *string `json:"phoneNumber"`
	Position                *string `json:"position"`
	StartWorkDate           *string `json:"start_work_date"`
	StartWorkDateCamel      *string `json:"startWorkDate"`
	SchoolDepartmentID      *string `json:"school_department_id"`
	SchoolDepartmentIDCamel *string `json:"schoolDepartmentId"`
	DepartmentID            *string `json:"department_id"`
	DepartmentIDCamel       *string `json:"departmentId"`
	IsActive                *bool   `json:"is_active"`
	IsActiveCamel           *bool   `json:"isActive"`
	RequestReason           *string `json:"request_reason"`
}

func firstNonEmptyUpdateString(values ...*string) *string {
	for _, value := range values {
		if value == nil {
			continue
		}
		trimmed := strings.TrimSpace(*value)
		if trimmed == "" {
			continue
		}
		v := trimmed
		return &v
	}
	return nil
}

func firstNonNilBool(values ...*bool) *bool {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
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

	req.MemberID = firstNonEmptyUpdateString(req.MemberID, req.MemberIDCamel)
	req.EmployeeCode = firstNonEmptyUpdateString(req.EmployeeCode, req.EmployeeCodeCamel)
	req.GenderID = firstNonEmptyUpdateString(req.GenderID, req.GenderIDCamel)
	req.PrefixID = firstNonEmptyUpdateString(req.PrefixID, req.PrefixIDCamel)
	req.FirstName = firstNonEmptyUpdateString(req.FirstName, req.FirstNameCamel)
	req.LastName = firstNonEmptyUpdateString(req.LastName, req.LastNameCamel)
	req.Phone = firstNonEmptyUpdateString(req.Phone, req.PhoneCamel)
	req.Position = firstNonEmptyUpdateString(req.Position)
	req.StartWorkDate = firstNonEmptyUpdateString(req.StartWorkDate, req.StartWorkDateCamel)
	req.SchoolDepartmentID = firstNonEmptyUpdateString(req.SchoolDepartmentID, req.SchoolDepartmentIDCamel)
	req.DepartmentID = firstNonEmptyUpdateString(req.DepartmentID, req.DepartmentIDCamel)
	req.IsActive = firstNonNilBool(req.IsActive, req.IsActiveCamel)

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	actorID := currentUser.Member.ID
	actorRole := currentUser.Member.Role
	switch actorRole {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin:
	default:
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, actorID, actorRole, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-management-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) MemberManagementsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

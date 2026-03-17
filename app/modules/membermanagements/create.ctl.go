package membermanagements

import (
	"strings"
	"time"

	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID                string  `json:"school_id"`
	SchoolIDCamel           string  `json:"schoolId"`
	Email                   string  `json:"email" binding:"required,email"`
	Password                string  `json:"password" binding:"required"`
	GenderID                string  `json:"gender_id"`
	GenderIDCamel           string  `json:"genderId"`
	PrefixID                string  `json:"prefix_id"`
	PrefixIDCamel           string  `json:"prefixId"`
	FirstName               *string `json:"first_name"`
	FirstNameCamel          *string `json:"firstName"`
	LastName                *string `json:"last_name"`
	LastNameCamel           *string `json:"lastName"`
	Phone                   *string `json:"phone"`
	PhoneCamel              *string `json:"phoneNumber"`
	Position                string  `json:"position" binding:"required"`
	StartWorkDate           string  `json:"start_work_date" binding:"required"`
	StartWorkDateCamel      string  `json:"startWorkDate"`
	SchoolDepartmentID      string  `json:"school_department_id"`
	SchoolDepartmentIDCamel string  `json:"schoolDepartmentId"`
	DepartmentID            string  `json:"department_id"`
	DepartmentIDCamel       string  `json:"departmentId"`
	RequestReason           *string `json:"request_reason"`
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func firstNonEmptyPointer(values ...*string) *string {
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

type registerManagementMemberResponse struct {
	ID        uuid.UUID      `json:"id"`
	SchoolID  uuid.UUID      `json:"school_id"`
	Email     string         `json:"email"`
	Role      ent.MemberRole `json:"role"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type registerManagementResponse struct {
	ID                 uuid.UUID  `json:"id"`
	MemberID           uuid.UUID  `json:"member_id"`
	EmployeeCode       string     `json:"employee_code"`
	GenderID           *uuid.UUID `json:"gender_id"`
	PrefixID           *uuid.UUID `json:"prefix_id"`
	FirstName          *string    `json:"first_name"`
	LastName           *string    `json:"last_name"`
	Phone              *string    `json:"phone"`
	Position           string     `json:"position"`
	StartWorkDate      time.Time  `json:"start_work_date"`
	SchoolDepartmentID uuid.UUID  `json:"school_department_id"`
	DepartmentID       uuid.UUID  `json:"department_id"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type registerManagementApprovalResponse struct {
	ID            uuid.UUID                 `json:"id"`
	CurrentStatus ent.ApprovalRequestStatus `json:"current_status"`
	SubmittedAt   *time.Time                `json:"submitted_at"`
	ResolvedAt    *time.Time                `json:"resolved_at"`
}

type registerManagementActionResponse struct {
	ID        uuid.UUID              `json:"id"`
	RequestID uuid.UUID              `json:"request_id"`
	Action    ent.ApprovalActionType `json:"action"`
	CreatedAt time.Time              `json:"created_at"`
}

type registerManagementResultResponse struct {
	Member          *registerManagementMemberResponse   `json:"member"`
	Management      *registerManagementResponse         `json:"management"`
	ApprovalRequest *registerManagementApprovalResponse `json:"approval_request,omitempty"`
	ApprovalAction  *registerManagementActionResponse   `json:"approval_action,omitempty"`
}

func toRegisterManagementResultResponse(item *ent.ManagementRegistrationResult) *registerManagementResultResponse {
	if item == nil || item.Member == nil || item.Management == nil {
		return nil
	}

	res := &registerManagementResultResponse{
		Member: &registerManagementMemberResponse{
			ID:        item.Member.ID,
			SchoolID:  item.Member.SchoolID,
			Email:     item.Member.Email,
			Role:      item.Member.Role,
			IsActive:  item.Member.IsActive,
			CreatedAt: item.Member.CreatedAt,
			UpdatedAt: item.Member.UpdatedAt,
		},
		Management: &registerManagementResponse{
			ID:                 item.Management.ID,
			MemberID:           item.Management.MemberID,
			EmployeeCode:       item.Management.EmployeeCode,
			GenderID:           item.Management.GenderID,
			PrefixID:           item.Management.PrefixID,
			FirstName:          item.Management.FirstName,
			LastName:           item.Management.LastName,
			Phone:              item.Management.Phone,
			Position:           item.Management.Position,
			StartWorkDate:      item.Management.StartWorkDate,
			SchoolDepartmentID: item.Management.SchoolDepartmentID,
			DepartmentID:       item.Management.DepartmentID,
			IsActive:           item.Management.IsActive,
			CreatedAt:          item.Management.CreatedAt,
			UpdatedAt:          item.Management.UpdatedAt,
		},
	}

	if item.Approval != nil {
		res.ApprovalRequest = &registerManagementApprovalResponse{
			ID:            item.Approval.ID,
			CurrentStatus: item.Approval.CurrentStatus,
			SubmittedAt:   item.Approval.SubmittedAt,
			ResolvedAt:    item.Approval.ResolvedAt,
		}
	}
	if item.ApprovalAction != nil {
		res.ApprovalAction = &registerManagementActionResponse{
			ID:        item.ApprovalAction.ID,
			RequestID: item.ApprovalAction.RequestID,
			Action:    item.ApprovalAction.Action,
			CreatedAt: item.ApprovalAction.CreatedAt,
		}
	}

	return res
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	actorID := currentUser.Member.ID
	actorRole := currentUser.Member.Role
	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	schoolIDRaw := firstNonEmptyString(req.SchoolID, req.SchoolIDCamel)
	schoolID := currentUser.Member.SchoolID
	if schoolIDRaw != "" {
		parsedSchoolID, err := uuid.Parse(schoolIDRaw)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = parsedSchoolID
	}
	if actorRole != ent.MemberRoleSuperadmin && schoolID != currentUser.Member.SchoolID {
		base.ValidateFailed(ctx, "invalid-school-scope", nil)
		return
	}
	schoolDepartmentIDRaw := firstNonEmptyString(req.SchoolDepartmentID, req.SchoolDepartmentIDCamel)
	var schoolDepartmentID *uuid.UUID
	if schoolDepartmentIDRaw != "" {
		parsed, err := uuid.Parse(schoolDepartmentIDRaw)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-department-id", nil)
			return
		}
		schoolDepartmentID = &parsed
	}

	departmentIDRaw := firstNonEmptyString(req.DepartmentID, req.DepartmentIDCamel)
	var departmentID *uuid.UUID
	if departmentIDRaw != "" {
		parsed, err := uuid.Parse(departmentIDRaw)
		if err != nil {
			base.BadRequest(ctx, "invalid-department-id", nil)
			return
		}
		departmentID = &parsed
	}

	genderIDRaw := firstNonEmptyString(req.GenderID, req.GenderIDCamel)
	var genderID *uuid.UUID
	if genderIDRaw != "" {
		parsed, err := uuid.Parse(genderIDRaw)
		if err != nil {
			base.BadRequest(ctx, "invalid-gender-id", nil)
			return
		}
		genderID = &parsed
	}

	prefixIDRaw := firstNonEmptyString(req.PrefixID, req.PrefixIDCamel)
	var prefixID *uuid.UUID
	if prefixIDRaw != "" {
		parsed, err := uuid.Parse(prefixIDRaw)
		if err != nil {
			base.BadRequest(ctx, "invalid-prefix-id", nil)
			return
		}
		prefixID = &parsed
	}

	firstName := firstNonEmptyPointer(req.FirstName, req.FirstNameCamel)
	lastName := firstNonEmptyPointer(req.LastName, req.LastNameCamel)
	phone := firstNonEmptyPointer(req.Phone, req.PhoneCamel)
	startWorkDate := firstNonEmptyString(req.StartWorkDate, req.StartWorkDateCamel)

	item, err := c.svc.Create(ctx.Request.Context(), actorID, actorRole, schoolID, req.Email, req.Password, genderID, prefixID, firstName, lastName, phone, req.Position, startWorkDate, schoolDepartmentID, departmentID, req.RequestReason)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-management-create-failed")
		return
	}

	base.Success(ctx, toRegisterManagementResultResponse(item), "success")
}

func (c *Controller) CreateMemberManagementController(ctx *gin.Context) {
	c.Create(ctx)
}

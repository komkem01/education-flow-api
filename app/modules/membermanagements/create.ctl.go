package membermanagements

import (
	"time"

	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	SchoolID      string  `json:"school_id" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	Password      string  `json:"password" binding:"required"`
	EmployeeCode  string  `json:"employee_code" binding:"required"`
	Position      string  `json:"position" binding:"required"`
	StartWorkDate string  `json:"start_work_date" binding:"required"`
	DepartmentID  string  `json:"department_id" binding:"required"`
	RequestReason *string `json:"request_reason"`
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
	ID            uuid.UUID `json:"id"`
	MemberID      uuid.UUID `json:"member_id"`
	EmployeeCode  string    `json:"employee_code"`
	Position      string    `json:"position"`
	StartWorkDate time.Time `json:"start_work_date"`
	DepartmentID  uuid.UUID `json:"department_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
			ID:            item.Management.ID,
			MemberID:      item.Management.MemberID,
			EmployeeCode:  item.Management.EmployeeCode,
			Position:      item.Management.Position,
			StartWorkDate: item.Management.StartWorkDate,
			DepartmentID:  item.Management.DepartmentID,
			IsActive:      item.Management.IsActive,
			CreatedAt:     item.Management.CreatedAt,
			UpdatedAt:     item.Management.UpdatedAt,
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

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, "invalid-school-id", nil)
		return
	}
	if actorRole != ent.MemberRoleSuperadmin && schoolID != currentUser.Member.SchoolID {
		base.ValidateFailed(ctx, "invalid-school-scope", nil)
		return
	}
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-department-id", nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), actorID, actorRole, schoolID, req.Email, req.Password, req.EmployeeCode, req.Position, req.StartWorkDate, departmentID, req.RequestReason)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-management-create-failed")
		return
	}

	base.Success(ctx, toRegisterManagementResultResponse(item), "success")
}

func (c *Controller) CreateMemberManagementController(ctx *gin.Context) {
	c.Create(ctx)
}

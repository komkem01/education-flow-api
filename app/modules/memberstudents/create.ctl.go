package memberstudents

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
	SchoolID              string  `json:"school_id" binding:"required"`
	Email                 string  `json:"email" binding:"required,email"`
	Password              string  `json:"password" binding:"required"`
	GenderID              string  `json:"gender_id" binding:"required"`
	PrefixID              string  `json:"prefix_id" binding:"required"`
	AdvisorTeacherID      string  `json:"advisor_teacher_id"`
	FirstNameTH           string  `json:"first_name_th" binding:"required"`
	LastNameTH            string  `json:"last_name_th" binding:"required"`
	FirstNameEN           *string `json:"first_name_en"`
	LastNameEN            *string `json:"last_name_en"`
	CitizenID             *string `json:"citizen_id"`
	Phone                 *string `json:"phone"`
	BirthDate             *string `json:"birth_date"`
	Nationality           *string `json:"nationality"`
	Religion              *string `json:"religion"`
	AddressCurrent        *string `json:"address_current"`
	AddressRegistered     *string `json:"address_registered"`
	EmergencyContactName  *string `json:"emergency_contact_name"`
	EmergencyContactPhone *string `json:"emergency_contact_phone"`
	ApprovalRequestReason *string `json:"approval_request_reason"`
}

type registerMemberResponse struct {
	ID        uuid.UUID      `json:"id"`
	SchoolID  uuid.UUID      `json:"school_id"`
	Email     string         `json:"email"`
	Role      ent.MemberRole `json:"role"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type registerStudentResponse struct {
	ID               uuid.UUID  `json:"id"`
	MemberID         uuid.UUID  `json:"member_id"`
	SchoolID         uuid.UUID  `json:"school_id"`
	GenderID         uuid.UUID  `json:"gender_id"`
	PrefixID         uuid.UUID  `json:"prefix_id"`
	AdvisorTeacherID *uuid.UUID `json:"advisor_teacher_id"`
	StudentCode      string     `json:"student_code"`
	FirstNameTH      string     `json:"first_name_th"`
	LastNameTH       string     `json:"last_name_th"`
	FirstNameEN      *string    `json:"first_name_en"`
	LastNameEN       *string    `json:"last_name_en"`
	CitizenID        *string    `json:"citizen_id"`
	Phone            *string    `json:"phone"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type registerStudentProfileResponse struct {
	ID                    uuid.UUID  `json:"id"`
	StudentID             uuid.UUID  `json:"student_id"`
	BirthDate             *time.Time `json:"birth_date"`
	Nationality           *string    `json:"nationality"`
	Religion              *string    `json:"religion"`
	AddressCurrent        *string    `json:"address_current"`
	AddressRegistered     *string    `json:"address_registered"`
	EmergencyContactName  *string    `json:"emergency_contact_name"`
	EmergencyContactPhone *string    `json:"emergency_contact_phone"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type registerApprovalResponse struct {
	ID            uuid.UUID                 `json:"id"`
	CurrentStatus ent.ApprovalRequestStatus `json:"current_status"`
	SubmittedAt   *time.Time                `json:"submitted_at"`
	ResolvedAt    *time.Time                `json:"resolved_at"`
}

type registerActionResponse struct {
	ID        uuid.UUID              `json:"id"`
	RequestID uuid.UUID              `json:"request_id"`
	Action    ent.ApprovalActionType `json:"action"`
	CreatedAt time.Time              `json:"created_at"`
}

type registerStudentResultResponse struct {
	Member          *registerMemberResponse         `json:"member"`
	Student         *registerStudentResponse        `json:"student"`
	Profile         *registerStudentProfileResponse `json:"profile"`
	ApprovalRequest *registerApprovalResponse       `json:"approval_request,omitempty"`
	ApprovalAction  *registerActionResponse         `json:"approval_action,omitempty"`
}

func toRegisterResultResponse(item *ent.StudentRegistrationResult) *registerStudentResultResponse {
	if item == nil || item.Member == nil || item.Student == nil || item.Profile == nil {
		return nil
	}

	res := &registerStudentResultResponse{
		Member: &registerMemberResponse{
			ID:        item.Member.ID,
			SchoolID:  item.Member.SchoolID,
			Email:     item.Member.Email,
			Role:      item.Member.Role,
			IsActive:  item.Member.IsActive,
			CreatedAt: item.Member.CreatedAt,
			UpdatedAt: item.Member.UpdatedAt,
		},
		Student: &registerStudentResponse{
			ID:               item.Student.ID,
			MemberID:         item.Student.MemberID,
			SchoolID:         item.Student.SchoolID,
			GenderID:         item.Student.GenderID,
			PrefixID:         item.Student.PrefixID,
			AdvisorTeacherID: item.Student.AdvisorTeacherID,
			StudentCode:      item.Student.StudentCode,
			FirstNameTH:      item.Student.FirstNameTH,
			LastNameTH:       item.Student.LastNameTH,
			FirstNameEN:      item.Student.FirstNameEN,
			LastNameEN:       item.Student.LastNameEN,
			CitizenID:        item.Student.CitizenID,
			Phone:            item.Student.Phone,
			IsActive:         item.Student.IsActive,
			CreatedAt:        item.Student.CreatedAt,
			UpdatedAt:        item.Student.UpdatedAt,
		},
		Profile: &registerStudentProfileResponse{
			ID:                    item.Profile.ID,
			StudentID:             item.Profile.StudentID,
			BirthDate:             item.Profile.BirthDate,
			Nationality:           item.Profile.Nationality,
			Religion:              item.Profile.Religion,
			AddressCurrent:        item.Profile.AddressCurrent,
			AddressRegistered:     item.Profile.AddressRegistered,
			EmergencyContactName:  item.Profile.EmergencyContactName,
			EmergencyContactPhone: item.Profile.EmergencyContactPhone,
			CreatedAt:             item.Profile.CreatedAt,
			UpdatedAt:             item.Profile.UpdatedAt,
		},
	}

	if item.Approval != nil {
		res.ApprovalRequest = &registerApprovalResponse{
			ID:            item.Approval.ID,
			CurrentStatus: item.Approval.CurrentStatus,
			SubmittedAt:   item.Approval.SubmittedAt,
			ResolvedAt:    item.Approval.ResolvedAt,
		}
	}
	if item.ApprovalAction != nil {
		res.ApprovalAction = &registerActionResponse{
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
	if !ok {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}
	if currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	actorID := currentUser.Member.ID
	actorRole := currentUser.Member.Role
	switch actorRole {
	case ent.MemberRoleSuperadmin, ent.MemberRoleAdmin, ent.MemberRoleTeacher:
	default:
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
	genderID, err := uuid.Parse(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, "invalid-gender-id", nil)
		return
	}
	prefixID, err := uuid.Parse(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, "invalid-prefix-id", nil)
		return
	}

	var advisorTeacherID *uuid.UUID
	if req.AdvisorTeacherID != "" {
		parsed, err := uuid.Parse(req.AdvisorTeacherID)
		if err != nil {
			base.BadRequest(ctx, "invalid-advisor-teacher-id", nil)
			return
		}
		advisorTeacherID = &parsed
	}

	item, err := c.svc.Create(
		ctx.Request.Context(),
		actorID,
		actorRole,
		schoolID,
		req.Email,
		req.Password,
		genderID,
		prefixID,
		advisorTeacherID,
		req.FirstNameTH,
		req.LastNameTH,
		req.FirstNameEN,
		req.LastNameEN,
		req.CitizenID,
		req.Phone,
		req.BirthDate,
		req.Nationality,
		req.Religion,
		req.AddressCurrent,
		req.AddressRegistered,
		req.EmergencyContactName,
		req.EmergencyContactPhone,
		req.ApprovalRequestReason,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-student-create-failed")
		return
	}

	base.Success(ctx, toRegisterResultResponse(item), "success")
}

func (c *Controller) CreateMemberStudentController(ctx *gin.Context) {
	c.Create(ctx)
}

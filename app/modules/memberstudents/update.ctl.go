package memberstudents

import (
	"eduflow/app/modules/auth"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	MemberID              *string `json:"member_id"`
	SchoolID              *string `json:"school_id"`
	GenderID              *string `json:"gender_id"`
	PrefixID              *string `json:"prefix_id"`
	AdvisorTeacherID      *string `json:"advisor_teacher_id"`
	StudentCode           *string `json:"student_code"`
	FirstNameTH           *string `json:"first_name_th"`
	LastNameTH            *string `json:"last_name_th"`
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
	BloodType             *string `json:"blood_type"`
	AllergyInfo           *string `json:"allergy_info"`
	ChronicDisease        *string `json:"chronic_disease"`
	MedicalNote           *string `json:"medical_note"`
	ApprovalRequestReason *string `json:"approval_request_reason"`
	IsActive              *bool   `json:"is_active"`
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

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser.Member == nil {
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

	item, err := c.svc.Update(ctx.Request.Context(), id, actorID, actorRole, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-student-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) MemberStudentsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

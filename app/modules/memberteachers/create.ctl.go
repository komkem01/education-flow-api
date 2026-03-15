package memberteachers

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
	SchoolID         string  `json:"school_id" binding:"required"`
	Email            string  `json:"email" binding:"required,email"`
	Password         string  `json:"password" binding:"required"`
	GenderID         string  `json:"gender_id" binding:"required"`
	PrefixID         string  `json:"prefix_id" binding:"required"`
	Code             string  `json:"code" binding:"required"`
	CitizenID        string  `json:"citizen_id" binding:"required"`
	FirstNameTH      string  `json:"first_name_th" binding:"required"`
	LastNameTH       string  `json:"last_name_th" binding:"required"`
	FirstNameEN      string  `json:"first_name_en" binding:"required"`
	LastNameEN       string  `json:"last_name_en" binding:"required"`
	Phone            string  `json:"phone" binding:"required"`
	Position         string  `json:"position" binding:"required"`
	AcademicStanding string  `json:"academic_standing" binding:"required"`
	DepartmentID     string  `json:"department_id" binding:"required"`
	StartDate        string  `json:"start_date" binding:"required"`
	EndDate          *string `json:"end_date"`
}

type registerTeacherMemberResponse struct {
	ID        uuid.UUID      `json:"id"`
	SchoolID  uuid.UUID      `json:"school_id"`
	Email     string         `json:"email"`
	Role      ent.MemberRole `json:"role"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type registerTeacherResponse struct {
	ID               uuid.UUID  `json:"id"`
	MemberID         uuid.UUID  `json:"member_id"`
	GenderID         uuid.UUID  `json:"gender_id"`
	PrefixID         uuid.UUID  `json:"prefix_id"`
	Code             string     `json:"code"`
	CitizenID        string     `json:"citizen_id"`
	FirstNameTH      string     `json:"first_name_th"`
	LastNameTH       string     `json:"last_name_th"`
	FirstNameEN      string     `json:"first_name_en"`
	LastNameEN       string     `json:"last_name_en"`
	Phone            string     `json:"phone"`
	Position         string     `json:"position"`
	AcademicStanding string     `json:"academic_standing"`
	DepartmentID     uuid.UUID  `json:"department_id"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type registerTeacherResultResponse struct {
	Member  *registerTeacherMemberResponse `json:"member"`
	Teacher *registerTeacherResponse       `json:"teacher"`
}

func toRegisterTeacherResultResponse(item *ent.TeacherRegistrationResult) *registerTeacherResultResponse {
	if item == nil || item.Member == nil || item.Teacher == nil {
		return nil
	}

	return &registerTeacherResultResponse{
		Member: &registerTeacherMemberResponse{
			ID:        item.Member.ID,
			SchoolID:  item.Member.SchoolID,
			Email:     item.Member.Email,
			Role:      item.Member.Role,
			IsActive:  item.Member.IsActive,
			CreatedAt: item.Member.CreatedAt,
			UpdatedAt: item.Member.UpdatedAt,
		},
		Teacher: &registerTeacherResponse{
			ID:               item.Teacher.ID,
			MemberID:         item.Teacher.MemberID,
			GenderID:         item.Teacher.GenderID,
			PrefixID:         item.Teacher.PrefixID,
			Code:             item.Teacher.Code,
			CitizenID:        item.Teacher.CitizenID,
			FirstNameTH:      item.Teacher.FirstNameTH,
			LastNameTH:       item.Teacher.LastNameTH,
			FirstNameEN:      item.Teacher.FirstNameEN,
			LastNameEN:       item.Teacher.LastNameEN,
			Phone:            item.Teacher.Phone,
			Position:         item.Teacher.Position,
			AcademicStanding: item.Teacher.AcademicStanding,
			DepartmentID:     item.Teacher.DepartmentID,
			StartDate:        item.Teacher.StartDate,
			EndDate:          item.Teacher.EndDate,
			IsActive:         item.Teacher.IsActive,
			CreatedAt:        item.Teacher.CreatedAt,
			UpdatedAt:        item.Teacher.UpdatedAt,
		},
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

	currentUser, ok := auth.CurrentUserFromGin(ctx)
	if !ok || currentUser.Member == nil {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

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
	departmentID, err := uuid.Parse(req.DepartmentID)
	if err != nil {
		base.BadRequest(ctx, "invalid-department-id", nil)
		return
	}

	teacher, err := c.svc.Create(
		ctx.Request.Context(),
		actorRole,
		schoolID,
		req.Email,
		req.Password,
		genderID,
		prefixID,
		req.Code,
		req.CitizenID,
		req.FirstNameTH,
		req.LastNameTH,
		req.FirstNameEN,
		req.LastNameEN,
		req.Phone,
		req.Position,
		req.AcademicStanding,
		departmentID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-teacher-create-failed")
		return
	}

	base.Success(ctx, toRegisterTeacherResultResponse(teacher), "success")
}

func (c *Controller) CreateMemberTeacherController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) MemberTeachersCreate(ctx *gin.Context) {
	c.Create(ctx)
}

package memberteachers

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
	MemberID              *uuid.UUID               `json:"member_id"`
	GenderID              *uuid.UUID               `json:"gender_id"`
	PrefixID              *uuid.UUID               `json:"prefix_id"`
	Code                  *string                  `json:"code"`
	CitizenID             *string                  `json:"citizen_id"`
	FirstNameTH           *string                  `json:"first_name_th"`
	LastNameTH            *string                  `json:"last_name_th"`
	FirstNameEN           *string                  `json:"first_name_en"`
	LastNameEN            *string                  `json:"last_name_en"`
	Phone                 *string                  `json:"phone"`
	Position              *string                  `json:"position"`
	AcademicStanding      *string                  `json:"academic_standing"`
	DepartmentID          *uuid.UUID               `json:"department_id"`
	StartDate             *string                  `json:"start_date"`
	EndDate               *string                  `json:"end_date"`
	IsActive              *bool                    `json:"is_active"`
	ApprovalRequestReason *string                  `json:"approval_request_reason"`
	Addresses             *[]TeacherAddressRequest `json:"addresses"`
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
	if actorRole != ent.MemberRoleAdmin && actorRole != ent.MemberRoleSuperadmin {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var addresses *[]ent.TeacherAddressInput
	if req.Addresses != nil {
		rows := make([]ent.TeacherAddressInput, 0, len(*req.Addresses))
		for _, addr := range *req.Addresses {
			rows = append(rows, ent.TeacherAddressInput{
				HouseNo:     addr.HouseNo,
				Village:     addr.Village,
				Road:        addr.Road,
				Province:    addr.Province,
				District:    addr.District,
				Subdistrict: addr.Subdistrict,
				PostalCode:  addr.PostalCode,
				IsPrimary:   addr.IsPrimary,
				SortOrder:   addr.SortOrder,
			})
		}
		addresses = &rows
	}

	teacher, err := c.svc.Update(ctx.Request.Context(), id, actorID, actorRole, &req, addresses)
	if err != nil {
		c.handleServiceError(ctx, log, err, "member-teacher-update-failed")
		return
	}

	base.Success(ctx, teacher, "success")
}

func (c *Controller) MemberTeachersUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

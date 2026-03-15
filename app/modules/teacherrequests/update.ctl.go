package teacherrequests

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateRequest struct {
	TeacherID     *string         `json:"teacher_id"`
	RequestType   *string         `json:"request_type"`
	RequestData   *map[string]any `json:"request_data"`
	RequestReason *string         `json:"request_reason"`
	Status        *string         `json:"status"`
	ApprovedBy    *string         `json:"approved_by"`
	ApprovedAt    *string         `json:"approved_at"`
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

	item, err := c.svc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		c.handleServiceError(ctx, log, err, "teacher-request-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) TeacherRequestsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

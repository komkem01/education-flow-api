package schoolannouncements

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type UpdateRequest struct {
	SchoolID      *string `json:"school_id"`
	Title         *string `json:"title"`
	Content       *string `json:"content"`
	Category      *string `json:"category"`
	Status        *string `json:"status"`
	AnnouncedAt   *string `json:"announced_at"`
	PublishedAt   *string `json:"published_at"`
	ExpiresAt     *string `json:"expires_at"`
	CreatedByName *string `json:"created_by_name"`
	TargetRole    *string `json:"target_role"`
	IsPinned      *bool   `json:"is_pinned"`
}

type UpdateByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func hasJSONNull(v json.RawMessage) bool {
	return strings.TrimSpace(string(v)) == "null"
}

func parseNullableDateField(raw map[string]json.RawMessage, key string) (*time.Time, bool, error) {
	v, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	if hasJSONNull(v) {
		return nil, true, nil
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return nil, false, err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, true, nil
	}
	t, err := parseDatePtr(&s)
	if err != nil {
		return nil, false, err
	}
	return t, false, nil
}

func parseNullableTargetRoleField(raw map[string]json.RawMessage) (*string, bool, error) {
	v, ok := raw["target_role"]
	if !ok {
		return nil, false, nil
	}
	if hasJSONNull(v) {
		return nil, true, nil
	}
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return nil, false, err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, true, nil
	}
	if !isValidTargetRole(s) {
		return nil, false, errors.New("invalid-target-role")
	}
	return &s, false, nil
}

func (c *Controller) Update(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	user, ok := c.currentUser(ctx)
	if !ok {
		return
	}
	if !c.canWrite(user) {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

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

	current, err := c.svc.Info(ctx.Request.Context(), id)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-announcement-info-failed")
		return
	}
	if user.Member.Role != ent.MemberRoleSuperadmin && current.SchoolID != user.Member.SchoolID {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	var req UpdateRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	rawPayload := map[string]json.RawMessage{}
	if err := ctx.ShouldBindBodyWith(&rawPayload, binding.JSON); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	var schoolID *uuid.UUID
	if req.SchoolID != nil {
		if *req.SchoolID == "" {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = &parsed
		if user.Member.Role != ent.MemberRoleSuperadmin && parsed != user.Member.SchoolID {
			base.Unauthorized(ctx, "unauthorized", nil)
			return
		}
	}

	var status *ent.SchoolAnnouncementStatus
	if req.Status != nil {
		if *req.Status == "" {
			base.BadRequest(ctx, "invalid-status", nil)
			return
		}
		parsed := ent.SchoolAnnouncementStatus(*req.Status)
		if !isValidStatus(parsed) {
			base.BadRequest(ctx, "invalid-status", nil)
			return
		}
		status = &parsed
	}

	announcedAt, clearAnnouncedAt, err := parseNullableDateField(rawPayload, "announced_at")
	if err != nil {
		base.BadRequest(ctx, "invalid-announced-at", nil)
		return
	}
	publishedAt, clearPublishedAt, err := parseNullableDateField(rawPayload, "published_at")
	if err != nil {
		base.BadRequest(ctx, "invalid-published-at", nil)
		return
	}
	expiresAt, clearExpiresAt, err := parseNullableDateField(rawPayload, "expires_at")
	if err != nil {
		base.BadRequest(ctx, "invalid-expires-at", nil)
		return
	}

	targetRole, clearTargetRole, err := parseNullableTargetRoleField(rawPayload)
	if err != nil {
		base.BadRequest(ctx, "invalid-target-role", nil)
		return
	}

	item, err := c.svc.Update(ctx.Request.Context(), id, &ent.SchoolAnnouncementUpdate{
		SchoolID:         schoolID,
		Title:            req.Title,
		Content:          req.Content,
		Category:         req.Category,
		Status:           status,
		AnnouncedAt:      announcedAt,
		ClearAnnouncedAt: clearAnnouncedAt,
		PublishedAt:      publishedAt,
		ClearPublishedAt: clearPublishedAt,
		ExpiresAt:        expiresAt,
		ClearExpiresAt:   clearExpiresAt,
		CreatedByName:    req.CreatedByName,
		TargetRole:       targetRole,
		ClearTargetRole:  clearTargetRole,
		IsPinned:         req.IsPinned,
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-announcement-update-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) UpdateSchoolAnnouncementController(ctx *gin.Context) {
	c.Update(ctx)
}

func (c *Controller) SchoolAnnouncementsUpdate(ctx *gin.Context) {
	c.Update(ctx)
}

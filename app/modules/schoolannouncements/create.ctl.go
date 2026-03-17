package schoolannouncements

import (
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
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

func (c *Controller) Create(ctx *gin.Context) {
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

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	schoolID := user.Member.SchoolID
	if req.SchoolID != nil && *req.SchoolID != "" {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			base.BadRequest(ctx, "invalid-school-id", nil)
			return
		}
		schoolID = parsed
	}

	if user.Member.Role != ent.MemberRoleSuperadmin && schoolID != user.Member.SchoolID {
		base.Unauthorized(ctx, "unauthorized", nil)
		return
	}

	status := ent.SchoolAnnouncementStatusDraft
	if req.Status != nil && *req.Status != "" {
		status = ent.SchoolAnnouncementStatus(*req.Status)
		if !isValidStatus(status) {
			base.BadRequest(ctx, "invalid-status", nil)
			return
		}
	}

	if req.TargetRole != nil && *req.TargetRole != "" && !isValidTargetRole(*req.TargetRole) {
		base.BadRequest(ctx, "invalid-target-role", nil)
		return
	}

	announcedAt, err := parseDatePtr(req.AnnouncedAt)
	if err != nil {
		base.BadRequest(ctx, "invalid-announced-at", nil)
		return
	}
	publishedAt, err := parseDatePtr(req.PublishedAt)
	if err != nil {
		base.BadRequest(ctx, "invalid-published-at", nil)
		return
	}
	expiresAt, err := parseDatePtr(req.ExpiresAt)
	if err != nil {
		base.BadRequest(ctx, "invalid-expires-at", nil)
		return
	}

	isPinned := false
	if req.IsPinned != nil {
		isPinned = *req.IsPinned
	}

	item, err := c.svc.Create(ctx.Request.Context(), &ent.SchoolAnnouncement{
		SchoolID:       schoolID,
		AuthorMemberID: user.Member.ID,
		Title:          req.Title,
		Content:        req.Content,
		Category:       req.Category,
		Status:         status,
		AnnouncedAt:    announcedAt,
		PublishedAt:    publishedAt,
		ExpiresAt:      expiresAt,
		CreatedByName:  req.CreatedByName,
		TargetRole:     req.TargetRole,
		IsPinned:       isPinned,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-announcement-create-failed")
		return
	}

	base.Success(ctx, item, "success")
}

func (c *Controller) CreateSchoolAnnouncementController(ctx *gin.Context) {
	c.Create(ctx)
}

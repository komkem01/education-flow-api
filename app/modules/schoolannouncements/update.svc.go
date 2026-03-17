package schoolannouncements

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, data *ent.SchoolAnnouncementUpdate) (*ent.SchoolAnnouncement, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "school_announcements.service.update")
	defer span.End()

	if data == nil {
		return nil, fmt.Errorf("%w", ErrSchoolAnnouncementInvalidUpdate)
	}

	if data.SchoolID == nil &&
		data.AuthorMemberID == nil &&
		data.Title == nil &&
		data.Content == nil &&
		data.Category == nil &&
		data.Status == nil &&
		data.AnnouncedAt == nil &&
		!data.ClearAnnouncedAt &&
		data.PublishedAt == nil &&
		!data.ClearPublishedAt &&
		data.ExpiresAt == nil &&
		!data.ClearExpiresAt &&
		data.CreatedByName == nil &&
		data.TargetRole == nil &&
		!data.ClearTargetRole &&
		data.IsPinned == nil {
		return nil, fmt.Errorf("%w", ErrSchoolAnnouncementInvalidUpdate)
	}

	item, err := s.db.UpdateSchoolAnnouncementByID(ctx, id, data)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

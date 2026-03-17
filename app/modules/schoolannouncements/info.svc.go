package schoolannouncements

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Info(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "school_announcements.service.info")
	defer span.End()

	item, err := s.db.GetSchoolAnnouncementByID(ctx, id)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

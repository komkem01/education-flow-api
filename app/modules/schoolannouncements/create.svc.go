package schoolannouncements

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
)

func (s *Service) Create(ctx context.Context, data *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "school_announcements.service.create")
	defer span.End()

	item, err := s.db.CreateSchoolAnnouncement(ctx, data)
	if err != nil {
		return nil, normalizeServiceError(err)
	}
	return item, nil
}

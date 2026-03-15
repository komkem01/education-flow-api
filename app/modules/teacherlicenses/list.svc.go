package teacherlicenses

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, status *string) ([]*ent.TeacherLicense, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherlicenses.service.list")
	defer span.End()

	var parsedStatus *ent.TeacherLicenseStatus
	if status != nil {
		v, ok := parseTeacherLicenseStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherLicenseConditionFail)
		}
		parsedStatus = &v
	}

	items, page, err := s.db.ListTeacherLicenses(ctx, req, teacherID, parsedStatus)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherLicenseConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

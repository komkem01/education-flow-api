package teacherrequests

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, req *base.RequestPaginate, teacherID *uuid.UUID, requestType *string, status *string) ([]*ent.TeacherRequest, *base.ResponsePaginate, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherrequests.service.list")
	defer span.End()

	var parsedType *ent.TeacherRequestType
	if requestType != nil {
		v, ok := parseTeacherRequestType(*requestType)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		parsedType = &v
	}

	var parsedStatus *ent.TeacherRequestStatus
	if status != nil {
		v, ok := parseTeacherRequestStatus(*status)
		if !ok {
			return nil, nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		parsedStatus = &v
	}

	items, page, err := s.db.ListTeacherRequests(ctx, req, teacherID, parsedType, parsedStatus)
	if err != nil {
		if base.IsPagErr(err) {
			return nil, nil, fmt.Errorf("%w: %v", ErrTeacherRequestConditionFail, err)
		}
		return nil, nil, normalizeServiceError(err)
	}

	return items, page, nil
}

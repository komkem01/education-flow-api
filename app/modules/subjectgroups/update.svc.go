package subjectgroups

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.SubjectGroup, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjectgroups.service.update")
	defer span.End()

	if req.SchoolID == nil && req.Code == nil && req.NameTH == nil && req.NameEN == nil && req.HeadTeacherID == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrSubjectGroupConditionFail)
	}

	payload := &ent.SubjectGroupUpdate{
		Code:     req.Code,
		NameTH:   req.NameTH,
		NameEN:   req.NameEN,
		IsActive: req.IsActive,
	}

	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrSubjectGroupConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if req.HeadTeacherID != nil {
		parsed, err := uuid.Parse(*req.HeadTeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrSubjectGroupConditionFail)
		}
		payload.HeadTeacherID = &parsed
	}

	item, err := s.db.UpdateSubjectGroupByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

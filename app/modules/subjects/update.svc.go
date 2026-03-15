package subjects

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.Subject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjects.service.update")
	defer span.End()

	if req.SchoolID == nil && req.SubjectGroupID == nil && req.Code == nil && req.NameTH == nil && req.NameEN == nil && req.Credit == nil && req.HoursPerWeek == nil && req.IsElective == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrSubjectConditionFail)
	}

	payload := &ent.SubjectUpdate{
		Code:         req.Code,
		NameTH:       req.NameTH,
		NameEN:       req.NameEN,
		Credit:       req.Credit,
		HoursPerWeek: req.HoursPerWeek,
		IsElective:   req.IsElective,
		IsActive:     req.IsActive,
	}

	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrSubjectConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if req.SubjectGroupID != nil {
		parsed, err := uuid.Parse(*req.SubjectGroupID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrSubjectConditionFail)
		}
		payload.SubjectGroupID = &parsed
	}

	item, err := s.db.UpdateSubjectByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

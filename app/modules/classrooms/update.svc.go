package classrooms

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest) (*ent.Classroom, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "classrooms.service.update")
	defer span.End()

	if req.SchoolID == nil && req.AcademicYearID == nil && req.Level == nil && req.RoomNo == nil && req.Name == nil && req.HomeroomTeacherID == nil && req.Capacity == nil && req.IsActive == nil {
		return nil, fmt.Errorf("%w", ErrClassroomConditionFail)
	}

	payload := &ent.ClassroomUpdate{
		Level:    req.Level,
		RoomNo:   req.RoomNo,
		Name:     req.Name,
		Capacity: req.Capacity,
		IsActive: req.IsActive,
	}

	if req.SchoolID != nil {
		parsed, err := uuid.Parse(*req.SchoolID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrClassroomConditionFail)
		}
		payload.SchoolID = &parsed
	}
	if req.AcademicYearID != nil {
		parsed, err := uuid.Parse(*req.AcademicYearID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrClassroomConditionFail)
		}
		payload.AcademicYearID = &parsed
	}
	if req.HomeroomTeacherID != nil {
		parsed, err := uuid.Parse(*req.HomeroomTeacherID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrClassroomConditionFail)
		}
		payload.HomeroomTeacherID = &parsed
	}

	item, err := s.db.UpdateClassroomByID(ctx, id, payload)
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

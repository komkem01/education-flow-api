package classrooms

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, academicYearID uuid.UUID, level string, roomNo string, name string, homeroomTeacherID *uuid.UUID, capacity *int, isActive bool) (*ent.Classroom, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "classrooms.service.create")
	defer span.End()

	item, err := s.db.CreateClassroom(ctx, &ent.Classroom{
		SchoolID:          schoolID,
		AcademicYearID:    academicYearID,
		Level:             level,
		RoomNo:            roomNo,
		Name:              name,
		HomeroomTeacherID: homeroomTeacherID,
		Capacity:          capacity,
		IsActive:          isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

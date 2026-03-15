package subjects

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, subjectGroupID uuid.UUID, code string, nameTH string, nameEN *string, credit float64, hoursPerWeek *int, isElective bool, isActive bool) (*ent.Subject, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjects.service.create")
	defer span.End()

	item, err := s.db.CreateSubject(ctx, &ent.Subject{
		SchoolID:       schoolID,
		SubjectGroupID: subjectGroupID,
		Code:           code,
		NameTH:         nameTH,
		NameEN:         nameEN,
		Credit:         credit,
		HoursPerWeek:   hoursPerWeek,
		IsElective:     isElective,
		IsActive:       isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

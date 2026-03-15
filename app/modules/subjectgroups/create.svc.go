package subjectgroups

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, schoolID uuid.UUID, code string, nameTH string, nameEN *string, headTeacherID *uuid.UUID, isActive bool) (*ent.SubjectGroup, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "subjectgroups.service.create")
	defer span.End()

	item, err := s.db.CreateSubjectGroup(ctx, &ent.SubjectGroup{
		SchoolID:      schoolID,
		Code:          code,
		NameTH:        nameTH,
		NameEN:        nameEN,
		HeadTeacherID: headTeacherID,
		IsActive:      isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

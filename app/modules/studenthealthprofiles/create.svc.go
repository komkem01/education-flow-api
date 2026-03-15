package studenthealthprofiles

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, studentID uuid.UUID, bloodType *string, allergyInfo *string, chronicDisease *string, medicalNote *string) (*ent.StudentHealthProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studenthealthprofiles.service.create")
	defer span.End()

	item, err := s.db.CreateStudentHealthProfile(ctx, &ent.StudentHealthProfile{
		StudentID:      studentID,
		BloodType:      bloodType,
		AllergyInfo:    allergyInfo,
		ChronicDisease: chronicDisease,
		MedicalNote:    medicalNote,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

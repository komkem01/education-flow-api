package studentprofiles

import (
	"context"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, studentID uuid.UUID, birthDate *string, nationality *string, religion *string, addressCurrent *string, addressRegistered *string, emergencyContactName *string, emergencyContactPhone *string) (*ent.StudentProfile, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "studentprofiles.service.create")
	defer span.End()

	var birthDateVal *time.Time
	if birthDate != nil && *birthDate != "" {
		parsed, err := time.Parse("2006-01-02", *birthDate)
		if err != nil {
			return nil, err
		}
		birthDateVal = &parsed
	}

	item, err := s.db.CreateStudentProfile(ctx, &ent.StudentProfile{
		StudentID:             studentID,
		BirthDate:             birthDateVal,
		Nationality:           nationality,
		Religion:              religion,
		AddressCurrent:        addressCurrent,
		AddressRegistered:     addressRegistered,
		EmergencyContactName:  emergencyContactName,
		EmergencyContactPhone: emergencyContactPhone,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

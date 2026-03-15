package memberstudents

import (
	"context"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, memberID uuid.UUID, schoolID uuid.UUID, genderID uuid.UUID, prefixID uuid.UUID, advisorTeacherID *uuid.UUID, studentCode string, firstNameTH string, lastNameTH string, firstNameEN *string, lastNameEN *string, citizenID *string, phone *string, isActive bool) (*ent.MemberStudent, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "memberstudents.service.create")
	defer span.End()

	item, err := s.db.CreateMemberStudent(ctx, &ent.MemberStudent{
		MemberID:         memberID,
		SchoolID:         schoolID,
		GenderID:         genderID,
		PrefixID:         prefixID,
		AdvisorTeacherID: advisorTeacherID,
		StudentCode:      studentCode,
		FirstNameTH:      firstNameTH,
		LastNameTH:       lastNameTH,
		FirstNameEN:      firstNameEN,
		LastNameEN:       lastNameEN,
		CitizenID:        citizenID,
		Phone:            phone,
		IsActive:         isActive,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

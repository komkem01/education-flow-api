package auth

import (
	"context"
	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils/base"
	"time"

	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Me(ctx context.Context, bearer string) (*MeResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "auth.service.me")
	defer span.End()

	user, err := s.resolveCurrentUser(ctx, bearer)
	if err != nil {
		return nil, err
	}

	res := &MeResponse{
		ID:        user.Member.ID.String(),
		SchoolID:  user.Member.SchoolID.String(),
		Email:     user.Member.Email,
		Role:      string(user.Member.Role),
		FullName:  displayNameFromEmail(user.Member.Email),
		IsActive:  user.Member.IsActive,
		LastLogin: user.Member.LastLogin,
		ExpireAt:  time.Unix(user.Claims.EXP, 0),
	}

	if s.school != nil {
		if school, err := s.school.GetSchoolByID(ctx, user.Member.SchoolID); err == nil && school != nil {
			res.SchoolName = school.Name
		}
	}

	if user.Member.Role == ent.MemberRoleTeacher && s.memberTeacher != nil {
		memberID := uuid.UUID(user.Member.ID)
		teachers, _, err := s.memberTeacher.ListMemberTeachers(ctx, &base.RequestPaginate{Page: 1, Size: 10}, nil, &memberID, nil)
		if err == nil && len(teachers) > 0 && teachers[0] != nil {
			res.FirstName = teachers[0].FirstNameTH
			res.LastName = teachers[0].LastNameTH
			res.FullName = teachers[0].FirstNameTH + " " + teachers[0].LastNameTH
		}
	}

	return res, nil
}

package approvals

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"

	"github.com/google/uuid"
)

func (s *Service) applyApprovalEffects(ctx context.Context, req *ent.ApprovalRequest) error {
	if req == nil {
		return nil
	}

	switch req.RequestType {
	case "student_registration":
		return s.activateStudentRegistration(ctx, req)
	case "management_registration":
		return s.activateManagementRegistration(ctx, req)
	default:
		return nil
	}
}

func (s *Service) activateStudentRegistration(ctx context.Context, req *ent.ApprovalRequest) error {
	memberID, err := payloadUUID(req.Payload, "member_id")
	if err != nil {
		return err
	}
	studentID, err := payloadUUID(req.Payload, "student_id")
	if err != nil {
		return err
	}

	active := true
	if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &active, nil); err != nil {
		return normalizeServiceError(err)
	}
	if _, err := s.studentDB.UpdateMemberStudentByID(ctx, studentID, &ent.MemberStudentUpdate{IsActive: &active}); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func (s *Service) activateManagementRegistration(ctx context.Context, req *ent.ApprovalRequest) error {
	memberID, err := payloadUUID(req.Payload, "member_id")
	if err != nil {
		return err
	}
	managementID, err := payloadUUID(req.Payload, "management_id")
	if err != nil {
		return err
	}

	active := true
	if _, err := s.memberDB.UpdateMemberByID(ctx, memberID, nil, nil, nil, nil, &active, nil); err != nil {
		return normalizeServiceError(err)
	}
	if _, err := s.managementDB.UpdateMemberManagementByID(ctx, managementID, &ent.MemberManagementUpdate{IsActive: &active}); err != nil {
		return normalizeServiceError(err)
	}

	return nil
}

func payloadUUID(payload map[string]any, key string) (uuid.UUID, error) {
	if payload == nil {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	raw, ok := payload[key]
	if !ok {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	v, ok := raw.(string)
	if !ok || v == "" {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w", ErrApprovalRequestConditionFail)
	}
	return id, nil
}

package teacherrequests

import (
	"context"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, teacherID uuid.UUID, requestedBy string, requestedByRole string, requestType string, requestData map[string]any, requestReason *string, status *string, approvedBy *uuid.UUID, approvedAt *string) (*ent.TeacherRequest, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "teacherrequests.service.create")
	defer span.End()

	requestedByID, err := uuid.Parse(requestedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}
	requestedRole, ok := parseApprovalActorRole(requestedByRole)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}

	parsedType, ok := parseTeacherRequestType(requestType)
	if !ok {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}

	parsedStatus := ent.TeacherRequestStatusPending
	if status != nil {
		v, ok := parseTeacherRequestStatus(*status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		if v != ent.TeacherRequestStatusPending {
			return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
		}
		parsedStatus = v
	}
	if approvedBy != nil || approvedAt != nil {
		return nil, fmt.Errorf("%w", ErrTeacherRequestConditionFail)
	}

	parsedApprovedAt, err := parseOptionalDateTime(approvedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTeacherRequestConditionFail, err)
	}

	item, err := s.db.CreateTeacherRequest(ctx, &ent.TeacherRequest{
		TeacherID:     teacherID,
		RequestType:   parsedType,
		RequestData:   requestData,
		RequestReason: requestReason,
		Status:        parsedStatus,
		ApprovedBy:    approvedBy,
		ApprovedAt:    parsedApprovedAt,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	now := time.Now()
	payload := map[string]any{
		"teacher_request_id": item.ID.String(),
		"teacher_id":         item.TeacherID.String(),
		"request_type":       item.RequestType,
		"request_data":       item.RequestData,
	}
	if item.RequestReason != nil {
		payload["request_reason"] = *item.RequestReason
	}

	approval, err := s.approval.CreateApprovalRequest(ctx, &ent.ApprovalRequest{
		RequestType:     "teacher_request",
		SubjectType:     "teacher_request",
		SubjectID:       &item.ID,
		RequestedBy:     requestedByID,
		RequestedByRole: requestedRole,
		Payload:         payload,
		CurrentStatus:   ent.ApprovalRequestStatusPending,
		SubmittedAt:     &now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	_, err = s.action.CreateApprovalAction(ctx, &ent.ApprovalAction{
		RequestID:   approval.ID,
		Action:      ent.ApprovalActionTypeSubmit,
		ActedBy:     requestedByID,
		ActedByRole: requestedRole,
		Comment:     requestReason,
		Metadata: map[string]any{
			"source":             "teacher_requests.create",
			"teacher_request_id": item.ID.String(),
		},
		CreatedAt: now,
	})
	if err != nil {
		return nil, normalizeServiceError(err)
	}

	return item, nil
}

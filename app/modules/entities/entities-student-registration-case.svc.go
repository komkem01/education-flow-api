package entities

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.StudentRegistrationCaseEntity = (*Service)(nil)

func (s *Service) CreateStudentRegistrationCase(ctx context.Context, data *ent.StudentRegistrationCase) (*ent.StudentRegistrationCase, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationCase, error) {
	row := new(ent.StudentRegistrationCase)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentRegistrationCases(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, studentID *uuid.UUID, status *ent.StudentRegistrationCaseStatus, registrationType *ent.StudentRegistrationType) ([]*ent.StudentRegistrationCase, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.StudentRegistrationCase, 0)
	query := s.db.NewSelect().Model(&items)

	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if studentID != nil {
		query.Where("student_id = ?", *studentID)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}
	if registrationType != nil {
		query.Where("registration_type = ?", *registrationType)
	}

	if err := req.SetSearchBy(query, []string{"case_no", "reason", "rejection_reason"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("requested_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"requested_at", "created_at", "case_no", "status", "registration_type"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID, data *ent.StudentRegistrationCaseUpdate) (*ent.StudentRegistrationCase, error) {
	query := s.db.NewUpdate().
		Model(&ent.StudentRegistrationCase{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.CaseNo != nil {
		query.Set("case_no = ?", *data.CaseNo)
	}
	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.StudentID != nil {
		query.Set("student_id = ?", *data.StudentID)
	}
	if data.RegistrationType != nil {
		query.Set("registration_type = ?", *data.RegistrationType)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.ApprovedBy != nil {
		query.Set("approved_by = ?", *data.ApprovedBy)
	}
	if data.RejectedBy != nil {
		query.Set("rejected_by = ?", *data.RejectedBy)
	}
	if data.SubmittedAt != nil {
		query.Set("submitted_at = ?", *data.SubmittedAt)
	}
	if data.ApprovedAt != nil {
		query.Set("approved_at = ?", *data.ApprovedAt)
	}
	if data.RejectedAt != nil {
		query.Set("rejected_at = ?", *data.RejectedAt)
	}
	if data.EffectiveDate != nil {
		query.Set("effective_date = ?", *data.EffectiveDate)
	}
	if data.Reason != nil {
		query.Set("reason = ?", *data.Reason)
	}
	if data.RejectionReason != nil {
		query.Set("rejection_reason = ?", *data.RejectionReason)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return s.GetStudentRegistrationCaseByID(ctx, id)
}

func (s *Service) SoftDeleteStudentRegistrationCaseByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.StudentRegistrationCase{}).
			Set("updated_at = now()").
			Where("id = ?", id).
			Exec(ctx)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return sql.ErrNoRows
		}

		_, err = tx.NewDelete().Model(&ent.StudentRegistrationCase{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

func (s *Service) GetStudentRegistrationCaseDetailByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationCaseDetail, error) {
	item, err := s.GetStudentRegistrationCaseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	detail := &ent.StudentRegistrationCaseDetail{Case: item}

	core := new(ent.StudentRegistrationStudentCore)
	if err := s.db.NewSelect().Model(core).Where("case_id = ?", id).Scan(ctx); err == nil {
		detail.Core = core
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	addresses := make([]*ent.StudentRegistrationAddress, 0)
	if err := s.db.NewSelect().Model(&addresses).Where("case_id = ?", id).Order("address_type ASC, sort_order ASC").Scan(ctx); err != nil {
		return nil, err
	}
	detail.Addresses = addresses

	health := new(ent.StudentRegistrationHealth)
	if err := s.db.NewSelect().Model(health).Where("case_id = ?", id).Scan(ctx); err == nil {
		detail.Health = health
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	guardians := make([]*ent.StudentRegistrationGuardian, 0)
	if err := s.db.NewSelect().Model(&guardians).Where("case_id = ?", id).Order("sort_order ASC").Scan(ctx); err != nil {
		return nil, err
	}
	detail.Guardians = guardians

	studentGuardians := make([]*ent.StudentRegistrationStudentGuardian, 0)
	if err := s.db.NewSelect().Model(&studentGuardians).Where("case_id = ?", id).Order("sort_order ASC").Scan(ctx); err != nil {
		return nil, err
	}
	detail.StudentGuardians = studentGuardians

	previousEducation := new(ent.StudentRegistrationPreviousEducation)
	if err := s.db.NewSelect().Model(previousEducation).Where("case_id = ?", id).Scan(ctx); err == nil {
		detail.PreviousEducation = previousEducation
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	familyEconomic := new(ent.StudentRegistrationFamilyEconomic)
	if err := s.db.NewSelect().Model(familyEconomic).Where("case_id = ?", id).Scan(ctx); err == nil {
		detail.FamilyEconomic = familyEconomic
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	documents := make([]*ent.StudentRegistrationDocument, 0)
	if err := s.db.NewSelect().Model(&documents).Where("case_id = ?", id).Order("created_at ASC").Scan(ctx); err != nil {
		return nil, err
	}
	detail.Documents = documents

	return detail, nil
}

func (s *Service) ReplaceStudentRegistrationCaseBundle(ctx context.Context, caseID uuid.UUID, bundle *ent.StudentRegistrationCaseBundle) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationStudentGuardian)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationGuardian)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationHealth)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationPreviousEducation)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationFamilyEconomic)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationDocument)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationAddress)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationStudentCore)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}

		guardianMap := map[uuid.UUID]uuid.UUID{}
		if bundle != nil && bundle.Core != nil {
			bundle.Core.ID = uuid.Nil
			bundle.Core.CaseID = caseID
			if _, err := tx.NewInsert().Model(bundle.Core).Exec(ctx); err != nil {
				return err
			}
		}

		if bundle != nil {
			for _, row := range bundle.Addresses {
				if row == nil {
					continue
				}
				row.ID = uuid.Nil
				row.CaseID = caseID
				if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
					return err
				}
			}

			if bundle.Health != nil {
				bundle.Health.ID = uuid.Nil
				bundle.Health.CaseID = caseID
				if _, err := tx.NewInsert().Model(bundle.Health).Exec(ctx); err != nil {
					return err
				}
			}

			if bundle.PreviousEducation != nil {
				bundle.PreviousEducation.ID = uuid.Nil
				bundle.PreviousEducation.CaseID = caseID
				if _, err := tx.NewInsert().Model(bundle.PreviousEducation).Exec(ctx); err != nil {
					return err
				}
			}

			if bundle.FamilyEconomic != nil {
				bundle.FamilyEconomic.ID = uuid.Nil
				bundle.FamilyEconomic.CaseID = caseID
				if _, err := tx.NewInsert().Model(bundle.FamilyEconomic).Exec(ctx); err != nil {
					return err
				}
			}

			for _, row := range bundle.Documents {
				if row == nil {
					continue
				}
				row.ID = uuid.Nil
				row.CaseID = caseID
				if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
					return err
				}
			}

			for _, row := range bundle.Guardians {
				if row == nil {
					continue
				}
				oldID := row.ID
				row.ID = uuid.Nil
				row.CaseID = caseID
				if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
					return err
				}
				if oldID != uuid.Nil {
					guardianMap[oldID] = row.ID
				}
			}

			for _, row := range bundle.StudentGuardians {
				if row == nil {
					continue
				}
				row.ID = uuid.Nil
				row.CaseID = caseID
				if mapped, ok := guardianMap[row.GuardianRowID]; ok {
					row.GuardianRowID = mapped
				}
				if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *Service) SubmitStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, reason *string) (*ent.StudentRegistrationCase, error) {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		item, err := s.getCaseForUpdate(ctx, tx, caseID)
		if err != nil {
			return err
		}
		if item.Status != ent.StudentRegistrationCaseStatusDraft && item.Status != ent.StudentRegistrationCaseStatusRejected {
			return fmt.Errorf("invalid-status")
		}

		now := time.Now()
		payload := map[string]any{
			"case_id":           item.ID.String(),
			"registration_type": item.RegistrationType,
			"school_id":         item.SchoolID.String(),
			"requested_by":      actorID.String(),
			"requested_by_role": actorRole,
		}
		if item.StudentID != nil {
			payload["student_id"] = item.StudentID.String()
		}
		if reason != nil && strings.TrimSpace(*reason) != "" {
			payload["reason"] = strings.TrimSpace(*reason)
		}

		approval := &ent.ApprovalRequest{
			RequestType:     "student_registration_case",
			SubjectType:     "student_registration_case",
			SubjectID:       &item.ID,
			RequestedBy:     actorID,
			RequestedByRole: actorRole,
			Payload:         payload,
			CurrentStatus:   ent.ApprovalRequestStatusPending,
			SubmittedAt:     &now,
		}
		if _, err := tx.NewInsert().Model(approval).Exec(ctx); err != nil {
			return err
		}

		action := &ent.ApprovalAction{
			RequestID:   approval.ID,
			Action:      ent.ApprovalActionTypeSubmit,
			ActedBy:     actorID,
			ActedByRole: actorRole,
			Comment:     reason,
			Metadata: map[string]any{
				"source":  "studentregistrationcases.submit",
				"case_id": item.ID.String(),
			},
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(action).Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationCase)(nil)).
			Set("status = ?", ent.StudentRegistrationCaseStatusPending).
			Set("submitted_at = ?", now).
			Set("updated_at = now()").
			Where("id = ?", item.ID).
			Exec(ctx); err != nil {
			return err
		}

		role := actorRole
		if _, err := tx.NewInsert().Model(&ent.StudentRegistrationAuditLog{
			CaseID:    item.ID,
			Action:    "submit",
			ActorID:   actorID,
			ActorRole: &role,
			OldStatus: &item.Status,
			NewStatus: ptrStudentCaseStatus(ent.StudentRegistrationCaseStatusPending),
			Comment:   reason,
			CreatedAt: now,
		}).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetStudentRegistrationCaseByID(ctx, caseID)
}

func (s *Service) RejectStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, reason *string) (*ent.StudentRegistrationCase, error) {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		item, err := s.getCaseForUpdate(ctx, tx, caseID)
		if err != nil {
			return err
		}
		if item.Status != ent.StudentRegistrationCaseStatusPending {
			return fmt.Errorf("invalid-status")
		}

		approval, err := s.getPendingCaseApproval(ctx, tx, caseID)
		if err != nil {
			return err
		}

		now := time.Now()
		if _, err := tx.NewUpdate().Model((*ent.ApprovalRequest)(nil)).
			Set("current_status = ?", ent.ApprovalRequestStatusRejected).
			Set("resolved_at = ?", now).
			Set("updated_at = now()").
			Where("id = ?", approval.ID).
			Exec(ctx); err != nil {
			return err
		}

		action := &ent.ApprovalAction{
			RequestID:   approval.ID,
			Action:      ent.ApprovalActionTypeReject,
			ActedBy:     actorID,
			ActedByRole: actorRole,
			Comment:     reason,
			Metadata: map[string]any{
				"source":  "studentregistrationcases.reject",
				"case_id": item.ID.String(),
			},
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(action).Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationCase)(nil)).
			Set("status = ?", ent.StudentRegistrationCaseStatusRejected).
			Set("rejected_by = ?", actorID).
			Set("rejected_at = ?", now).
			Set("rejection_reason = ?", reason).
			Set("updated_at = now()").
			Where("id = ?", item.ID).
			Exec(ctx); err != nil {
			return err
		}

		role := actorRole
		if _, err := tx.NewInsert().Model(&ent.StudentRegistrationAuditLog{
			CaseID:    item.ID,
			Action:    "reject",
			ActorID:   actorID,
			ActorRole: &role,
			OldStatus: &item.Status,
			NewStatus: ptrStudentCaseStatus(ent.StudentRegistrationCaseStatusRejected),
			Comment:   reason,
			CreatedAt: now,
		}).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetStudentRegistrationCaseByID(ctx, caseID)
}

func (s *Service) ApproveAndApplyStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string, idempotencyKey *string, metadata map[string]any) (*ent.StudentRegistrationCase, error) {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		item, err := s.getCaseForUpdate(ctx, tx, caseID)
		if err != nil {
			return err
		}
		if item.Status != ent.StudentRegistrationCaseStatusPending {
			return fmt.Errorf("invalid-status")
		}

		approval, err := s.getPendingCaseApproval(ctx, tx, caseID)
		if err != nil {
			return err
		}

		if idempotencyKey != nil && strings.TrimSpace(*idempotencyKey) != "" {
			exists, checkErr := s.hasApprovalIdempotencyAction(ctx, tx, approval.ID, ent.ApprovalActionTypeApprove, strings.TrimSpace(*idempotencyKey))
			if checkErr != nil {
				return checkErr
			}
			if exists {
				return nil
			}
		}

		now := time.Now()
		if _, err := tx.NewUpdate().Model((*ent.ApprovalRequest)(nil)).
			Set("current_status = ?", ent.ApprovalRequestStatusApproved).
			Set("resolved_at = ?", now).
			Set("updated_at = now()").
			Where("id = ?", approval.ID).
			Exec(ctx); err != nil {
			return err
		}

		actionMetadata := map[string]any{"source": "studentregistrationcases.approve", "case_id": item.ID.String()}
		for k, v := range metadata {
			actionMetadata[k] = v
		}
		if idempotencyKey != nil && strings.TrimSpace(*idempotencyKey) != "" {
			actionMetadata["idempotency_key"] = strings.TrimSpace(*idempotencyKey)
		}

		if _, err := tx.NewInsert().Model(&ent.ApprovalAction{
			RequestID:   approval.ID,
			Action:      ent.ApprovalActionTypeApprove,
			ActedBy:     actorID,
			ActedByRole: actorRole,
			Comment:     comment,
			Metadata:    actionMetadata,
			CreatedAt:   now,
		}).Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationCase)(nil)).
			Set("status = ?", ent.StudentRegistrationCaseStatusApproved).
			Set("approved_by = ?", actorID).
			Set("approved_at = ?", now).
			Set("updated_at = now()").
			Where("id = ?", item.ID).
			Exec(ctx); err != nil {
			return err
		}

		role := actorRole
		if _, err := tx.NewInsert().Model(&ent.StudentRegistrationAuditLog{
			CaseID:    item.ID,
			Action:    "approve",
			ActorID:   actorID,
			ActorRole: &role,
			OldStatus: &item.Status,
			NewStatus: ptrStudentCaseStatus(ent.StudentRegistrationCaseStatusApproved),
			Comment:   comment,
			CreatedAt: now,
		}).Exec(ctx); err != nil {
			return err
		}

		return s.applyStudentRegistrationCaseInTx(ctx, tx, item.ID, actorID, actorRole, comment)
	})
	if err != nil {
		return nil, err
	}
	return s.GetStudentRegistrationCaseByID(ctx, caseID)
}

func (s *Service) ApplyStudentRegistrationCase(ctx context.Context, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string) (*ent.StudentRegistrationCase, error) {
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		item, err := s.getCaseForUpdate(ctx, tx, caseID)
		if err != nil {
			return err
		}
		if item.Status != ent.StudentRegistrationCaseStatusApproved {
			return fmt.Errorf("invalid-status")
		}
		return s.applyStudentRegistrationCaseInTx(ctx, tx, caseID, actorID, actorRole, comment)
	})
	if err != nil {
		return nil, err
	}
	return s.GetStudentRegistrationCaseByID(ctx, caseID)
}

func (s *Service) applyStudentRegistrationCaseInTx(ctx context.Context, tx bun.Tx, caseID uuid.UUID, actorID uuid.UUID, actorRole ent.ApprovalActorRole, comment *string) error {
	item := new(ent.StudentRegistrationCase)
	if err := tx.NewSelect().Model(item).Where("id = ?", caseID).Scan(ctx); err != nil {
		return err
	}

	core := new(ent.StudentRegistrationStudentCore)
	if err := tx.NewSelect().Model(core).Where("case_id = ?", caseID).Scan(ctx); err != nil {
		return err
	}

	studentID := item.StudentID
	if studentID == nil {
		studentID = core.StudentID
	}
	if studentID == nil && item.RegistrationType == ent.StudentRegistrationTypeNewEnrollment {
		if core.PendingMemberEmail == nil || core.PendingMemberPasswordHash == nil {
			return fmt.Errorf("missing-student-credentials")
		}

		email := strings.ToLower(strings.TrimSpace(*core.PendingMemberEmail))
		if email == "" {
			return fmt.Errorf("missing-student-credentials")
		}

		member := &ent.Member{
			SchoolID: item.SchoolID,
			Email:    email,
			Password: *core.PendingMemberPasswordHash,
			Role:     ent.MemberRoleStudent,
			IsActive: core.IsActiveTarget,
		}
		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return err
		}

		const maxStudentCodeRetry = 10
		var createdStudent *ent.MemberStudent
		for i := 0; i < maxStudentCodeRetry; i += 1 {
			studentCode, codeErr := utils.GenerateNumericCode("STD", 6)
			if codeErr != nil {
				return codeErr
			}

			candidate := &ent.MemberStudent{
				MemberID:         member.ID,
				SchoolID:         item.SchoolID,
				GenderID:         core.GenderID,
				PrefixID:         core.PrefixID,
				AdvisorTeacherID: core.AdvisorTeacherID,
				StudentCode:      studentCode,
				FirstNameTH:      core.FirstNameTH,
				LastNameTH:       core.LastNameTH,
				FirstNameEN:      core.FirstNameEN,
				LastNameEN:       core.LastNameEN,
				CitizenID:        core.CitizenID,
				Phone:            core.Phone,
				IsActive:         core.IsActiveTarget,
			}
			if _, err := tx.NewInsert().Model(candidate).Exec(ctx); err != nil {
				errText := strings.ToLower(err.Error())
				if strings.Contains(errText, "student_code") && strings.Contains(errText, "duplicate") {
					continue
				}
				return err
			}
			createdStudent = candidate
			break
		}
		if createdStudent == nil {
			return fmt.Errorf("duplicate-student-code")
		}

		studentID = &createdStudent.ID
		core.StudentID = studentID
		core.MemberID = &member.ID

		if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationCase)(nil)).
			Set("student_id = ?", *studentID).
			Set("updated_at = now()").
			Where("id = ?", caseID).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationStudentCore)(nil)).
			Set("student_id = ?", *studentID).
			Set("member_id = ?", member.ID).
			Set("updated_at = now()").
			Where("id = ?", core.ID).
			Exec(ctx); err != nil {
			return err
		}
	}

	if studentID == nil {
		return fmt.Errorf("missing-student-id")
	}

	targetStudent := new(ent.MemberStudent)
	if err := tx.NewSelect().Model(targetStudent).Where("id = ?", *studentID).Limit(1).Scan(ctx); err != nil {
		return err
	}

	res, err := tx.NewUpdate().Model((*ent.MemberStudent)(nil)).
		Set("gender_id = ?", core.GenderID).
		Set("prefix_id = ?", core.PrefixID).
		Set("advisor_teacher_id = ?", core.AdvisorTeacherID).
		Set("first_name_th = ?", core.FirstNameTH).
		Set("last_name_th = ?", core.LastNameTH).
		Set("first_name_en = ?", core.FirstNameEN).
		Set("last_name_en = ?", core.LastNameEN).
		Set("citizen_id = ?", core.CitizenID).
		Set("phone = ?", core.Phone).
		Set("is_active = ?", core.IsActiveTarget).
		Set("updated_at = now()").
		Where("id = ?", *studentID).
		Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	memberID := core.MemberID
	if memberID == nil {
		memberID = &targetStudent.MemberID
	}
	if memberID != nil {
		query := tx.NewUpdate().Model((*ent.Member)(nil)).
			Set("is_active = ?", core.IsActiveTarget).
			Set("updated_at = now()").
			Where("id = ?", *memberID)

		if core.PendingMemberEmail != nil {
			query.Set("email = ?", strings.ToLower(strings.TrimSpace(*core.PendingMemberEmail)))
		}
		if core.PendingMemberPasswordHash != nil {
			query.Set("password = ?", *core.PendingMemberPasswordHash)
		}

		if _, err := query.Exec(ctx); err != nil {
			return err
		}
	}

	currentAddress, _ := s.composeRegistrationAddressText(ctx, tx, caseID, ent.RegistrationAddressTypeCurrent)
	registeredAddress, _ := s.composeRegistrationAddressText(ctx, tx, caseID, ent.RegistrationAddressTypeRegistered)

	profile := new(ent.StudentProfile)
	profileErr := tx.NewSelect().Model(profile).Where("student_id = ?", *studentID).Limit(1).Scan(ctx)
	if profileErr == nil {
		if _, err := tx.NewUpdate().Model((*ent.StudentProfile)(nil)).
			Set("address_current = ?", currentAddress).
			Set("address_registered = ?", registeredAddress).
			Set("updated_at = now()").
			Where("id = ?", profile.ID).
			Exec(ctx); err != nil {
			return err
		}
	} else if profileErr == sql.ErrNoRows {
		if _, err := tx.NewInsert().Model(&ent.StudentProfile{
			StudentID:         *studentID,
			AddressCurrent:    currentAddress,
			AddressRegistered: registeredAddress,
		}).Exec(ctx); err != nil {
			return err
		}
	} else {
		return profileErr
	}

	health := new(ent.StudentRegistrationHealth)
	if err := tx.NewSelect().Model(health).Where("case_id = ?", caseID).Scan(ctx); err == nil {
		targetHealth := new(ent.StudentHealthProfile)
		healthErr := tx.NewSelect().Model(targetHealth).Where("student_id = ?", *studentID).Limit(1).Scan(ctx)
		if healthErr == nil {
			if _, err := tx.NewUpdate().Model((*ent.StudentHealthProfile)(nil)).
				Set("blood_type = ?", health.BloodType).
				Set("allergy_info = ?", health.AllergyInfo).
				Set("chronic_disease = ?", health.ChronicDisease).
				Set("medical_note = ?", health.MedicalNote).
				Set("updated_at = now()").
				Where("id = ?", targetHealth.ID).
				Exec(ctx); err != nil {
				return err
			}
		} else if healthErr == sql.ErrNoRows {
			if _, err := tx.NewInsert().Model(&ent.StudentHealthProfile{
				StudentID:      *studentID,
				BloodType:      health.BloodType,
				AllergyInfo:    health.AllergyInfo,
				ChronicDisease: health.ChronicDisease,
				MedicalNote:    health.MedicalNote,
			}).Exec(ctx); err != nil {
				return err
			}
		} else {
			return healthErr
		}
	} else if err != sql.ErrNoRows {
		return err
	}

	if _, err := tx.NewDelete().Model((*ent.StudentGuardian)(nil)).Where("student_id = ?", *studentID).Exec(ctx); err != nil {
		return err
	}

	guardians := make([]*ent.StudentRegistrationGuardian, 0)
	if err := tx.NewSelect().Model(&guardians).Where("case_id = ?", caseID).Order("sort_order ASC").Scan(ctx); err != nil {
		return err
	}
	links := make([]*ent.StudentRegistrationStudentGuardian, 0)
	if err := tx.NewSelect().Model(&links).Where("case_id = ?", caseID).Order("sort_order ASC").Scan(ctx); err != nil {
		return err
	}

	guardianMap := make(map[uuid.UUID]uuid.UUID, len(guardians))
	for _, row := range guardians {
		if row == nil {
			continue
		}
		phone := row.Phone
		created := &ent.MemberGuardian{
			SchoolID:    item.SchoolID,
			GenderID:    row.GenderID,
			PrefixID:    row.PrefixID,
			FirstNameTH: row.FirstNameTH,
			LastNameTH:  row.LastNameTH,
			FirstNameEN: row.FirstNameEN,
			LastNameEN:  row.LastNameEN,
			CitizenID:   row.CitizenID,
			Phone:       &phone,
			IsActive:    row.IsActiveTarget,
		}
		if _, err := tx.NewInsert().Model(created).Exec(ctx); err != nil {
			return err
		}
		guardianMap[row.ID] = created.ID
	}

	for _, row := range links {
		if row == nil {
			continue
		}
		newGuardianID, ok := guardianMap[row.GuardianRowID]
		if !ok {
			continue
		}
		if _, err := tx.NewInsert().Model(&ent.StudentGuardian{
			StudentID:          *studentID,
			GuardianID:         newGuardianID,
			Relationship:       row.Relationship,
			IsMainGuardian:     row.IsMainGuardian,
			CanPickup:          row.CanPickup,
			IsEmergencyContact: row.IsEmergencyContact,
			Note:               row.Note,
		}).Exec(ctx); err != nil {
			return err
		}
	}

	oldStatus := item.Status
	newStatus := ent.StudentRegistrationCaseStatusApplied
	now := time.Now()
	if _, err := tx.NewUpdate().Model((*ent.StudentRegistrationCase)(nil)).
		Set("status = ?", newStatus).
		Set("updated_at = now()").
		Where("id = ?", caseID).
		Exec(ctx); err != nil {
		return err
	}

	role := actorRole
	if _, err := tx.NewInsert().Model(&ent.StudentRegistrationAuditLog{
		CaseID:    caseID,
		Action:    "apply",
		ActorID:   actorID,
		ActorRole: &role,
		OldStatus: &oldStatus,
		NewStatus: &newStatus,
		Comment:   comment,
		CreatedAt: now,
	}).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) composeRegistrationAddressText(ctx context.Context, tx bun.Tx, caseID uuid.UUID, addressType ent.RegistrationAddressType) (*string, error) {
	row := new(ent.StudentRegistrationAddress)
	err := tx.NewSelect().Model(row).
		Where("case_id = ?", caseID).
		Where("address_type = ?", addressType).
		Order("is_primary DESC, sort_order ASC").
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parts := []string{strings.TrimSpace(row.HouseNo)}
	if row.Village != nil && strings.TrimSpace(*row.Village) != "" {
		parts = append(parts, strings.TrimSpace(*row.Village))
	}
	if row.Road != nil && strings.TrimSpace(*row.Road) != "" {
		parts = append(parts, strings.TrimSpace(*row.Road))
	}
	parts = append(parts, strings.TrimSpace(row.Subdistrict), strings.TrimSpace(row.District), strings.TrimSpace(row.Province), strings.TrimSpace(row.PostalCode))
	joined := strings.Join(parts, ", ")
	return &joined, nil
}

func (s *Service) getCaseForUpdate(ctx context.Context, tx bun.Tx, caseID uuid.UUID) (*ent.StudentRegistrationCase, error) {
	item := new(ent.StudentRegistrationCase)
	if err := tx.NewSelect().Model(item).Where("id = ?", caseID).For("UPDATE").Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) getPendingCaseApproval(ctx context.Context, tx bun.Tx, caseID uuid.UUID) (*ent.ApprovalRequest, error) {
	item := new(ent.ApprovalRequest)
	if err := tx.NewSelect().Model(item).
		Where("subject_type = ?", "student_registration_case").
		Where("subject_id = ?", caseID).
		Where("current_status = ?", ent.ApprovalRequestStatusPending).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) hasApprovalIdempotencyAction(ctx context.Context, tx bun.Tx, requestID uuid.UUID, action ent.ApprovalActionType, key string) (bool, error) {
	count, err := tx.NewSelect().Model((*ent.ApprovalAction)(nil)).
		Where("request_id = ?", requestID).
		Where("action = ?", action).
		Where("metadata->>'idempotency_key' = ?", key).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ptrStudentCaseStatus(v ent.StudentRegistrationCaseStatus) *ent.StudentRegistrationCaseStatus {
	return &v
}

func (s *Service) GetStudentRegistrationPreviousEducationByCaseID(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationPreviousEducation, error) {
	row := new(ent.StudentRegistrationPreviousEducation)
	if err := s.db.NewSelect().Model(row).Where("case_id = ?", caseID).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpsertStudentRegistrationPreviousEducationByCaseID(ctx context.Context, caseID uuid.UUID, data *ent.StudentRegistrationPreviousEducation) (*ent.StudentRegistrationPreviousEducation, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid-data")
	}

	current, err := s.GetStudentRegistrationPreviousEducationByCaseID(ctx, caseID)
	if err == nil {
		_, err = s.db.NewUpdate().
			Model((*ent.StudentRegistrationPreviousEducation)(nil)).
			Set("previous_school_name = ?", data.PreviousSchoolName).
			Set("previous_school_province = ?", data.PreviousSchoolProvince).
			Set("previous_grade_level = ?", data.PreviousGradeLevel).
			Set("gpa = ?", data.GPA).
			Set("transfer_certificate_no = ?", data.TransferCertificateNo).
			Set("transfer_date = ?", data.TransferDate).
			Set("transcript_received = ?", data.TranscriptReceived).
			Set("remarks = ?", data.Remarks).
			Set("updated_at = now()").
			Where("id = ?", current.ID).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return s.GetStudentRegistrationPreviousEducationByCaseID(ctx, caseID)
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	data.ID = uuid.Nil
	data.CaseID = caseID
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentRegistrationFamilyEconomicByCaseID(ctx context.Context, caseID uuid.UUID) (*ent.StudentRegistrationFamilyEconomic, error) {
	row := new(ent.StudentRegistrationFamilyEconomic)
	if err := s.db.NewSelect().Model(row).Where("case_id = ?", caseID).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpsertStudentRegistrationFamilyEconomicByCaseID(ctx context.Context, caseID uuid.UUID, data *ent.StudentRegistrationFamilyEconomic) (*ent.StudentRegistrationFamilyEconomic, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid-data")
	}

	current, err := s.GetStudentRegistrationFamilyEconomicByCaseID(ctx, caseID)
	if err == nil {
		_, err = s.db.NewUpdate().
			Model((*ent.StudentRegistrationFamilyEconomic)(nil)).
			Set("household_size = ?", data.HouseholdSize).
			Set("household_income_monthly = ?", data.HouseholdIncomeMonthly).
			Set("income_bracket = ?", data.IncomeBracket).
			Set("scholarship_flag = ?", data.ScholarshipFlag).
			Set("scholarship_type = ?", data.ScholarshipType).
			Set("welfare_flag = ?", data.WelfareFlag).
			Set("welfare_type = ?", data.WelfareType).
			Set("debt_flag = ?", data.DebtFlag).
			Set("debt_detail = ?", data.DebtDetail).
			Set("updated_at = now()").
			Where("id = ?", current.ID).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
		return s.GetStudentRegistrationFamilyEconomicByCaseID(ctx, caseID)
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	data.ID = uuid.Nil
	data.CaseID = caseID
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) ListStudentRegistrationDocumentsByCaseID(ctx context.Context, caseID uuid.UUID) ([]*ent.StudentRegistrationDocument, error) {
	rows := make([]*ent.StudentRegistrationDocument, 0)
	if err := s.db.NewSelect().Model(&rows).Where("case_id = ?", caseID).Order("created_at ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *Service) ReplaceStudentRegistrationDocumentsByCaseID(ctx context.Context, caseID uuid.UUID, rows []*ent.StudentRegistrationDocument) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().Model((*ent.StudentRegistrationDocument)(nil)).Where("case_id = ?", caseID).Exec(ctx); err != nil {
			return err
		}
		for _, row := range rows {
			if row == nil {
				continue
			}
			row.ID = uuid.Nil
			row.CaseID = caseID
			if _, err := tx.NewInsert().Model(row).Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) CreateStudentRegistrationRule(ctx context.Context, data *ent.StudentRegistrationRule) (*ent.StudentRegistrationRule, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID) (*ent.StudentRegistrationRule, error) {
	row := new(ent.StudentRegistrationRule)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentRegistrationRules(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, registrationType *ent.StudentRegistrationType) ([]*ent.StudentRegistrationRule, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	rows := make([]*ent.StudentRegistrationRule, 0)
	query := s.db.NewSelect().Model(&rows)
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if registrationType != nil {
		query.Where("registration_type = ?", *registrationType)
	}

	if err := req.SetSearchBy(query, []string{"field_code", "validation_message"}); err != nil {
		return nil, nil, err
	}
	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "active_from", "field_code", "registration_type"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return rows, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID, data *ent.StudentRegistrationRuleUpdate) (*ent.StudentRegistrationRule, error) {
	query := s.db.NewUpdate().
		Model((*ent.StudentRegistrationRule)(nil)).
		Set("updated_at = now()").
		Where("id = ?", id)

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.RegistrationType != nil {
		query.Set("registration_type = ?", *data.RegistrationType)
	}
	if data.FieldCode != nil {
		query.Set("field_code = ?", *data.FieldCode)
	}
	if data.IsRequired != nil {
		query.Set("is_required = ?", *data.IsRequired)
	}
	if data.ValidationRegex != nil {
		query.Set("validation_regex = ?", *data.ValidationRegex)
	}
	if data.ValidationMessage != nil {
		query.Set("validation_message = ?", *data.ValidationMessage)
	}
	if data.ActiveFrom != nil {
		query.Set("active_from = ?", *data.ActiveFrom)
	}
	if data.ActiveTo != nil {
		query.Set("active_to = ?", *data.ActiveTo)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return s.GetStudentRegistrationRuleByID(ctx, id)
}

func (s *Service) SoftDeleteStudentRegistrationRuleByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model((*ent.StudentRegistrationRule)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

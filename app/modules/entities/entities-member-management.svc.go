package entities

import (
	"context"
	"database/sql"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberManagementEntity = (*Service)(nil)

func (s *Service) CreateMemberManagement(ctx context.Context, data *ent.MemberManagement) (*ent.MemberManagement, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) RegisterManagement(ctx context.Context, data *ent.ManagementRegistrationInput) (*ent.ManagementRegistrationResult, error) {
	result := new(ent.ManagementRegistrationResult)

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := &ent.Member{
			SchoolID:  data.MemberSchoolID,
			Email:     data.MemberEmail,
			Password:  data.MemberPasswordHash,
			Role:      data.MemberRole,
			IsActive:  data.MemberIsActive,
			LastLogin: data.MemberLastLogin,
		}
		if _, err := tx.NewInsert().Model(member).Returning("*").Exec(ctx); err != nil {
			return err
		}

		management := &ent.MemberManagement{
			MemberID:      member.ID,
			EmployeeCode:  data.ManagementEmployeeCode,
			Position:      data.ManagementPosition,
			StartWorkDate: data.ManagementStartWorkDate,
			DepartmentID:  data.ManagementDepartmentID,
			IsActive:      data.ManagementIsActive,
		}
		if _, err := tx.NewInsert().Model(management).Returning("*").Exec(ctx); err != nil {
			return err
		}

		now := time.Now()
		subjectID := management.ID
		approval := &ent.ApprovalRequest{
			RequestType:     "management_registration",
			SubjectType:     "member_management",
			SubjectID:       &subjectID,
			RequestedBy:     data.RequestedBy,
			RequestedByRole: data.RequestedByRole,
			Payload: map[string]any{
				"member_id":     member.ID.String(),
				"management_id": management.ID.String(),
				"employee_code": management.EmployeeCode,
				"department_id": management.DepartmentID.String(),
				"target_role":   string(member.Role),
			},
			CurrentStatus: ent.ApprovalRequestStatusPending,
			SubmittedAt:   &now,
		}
		if data.RequestReason != nil && *data.RequestReason != "" {
			approval.Payload["reason"] = *data.RequestReason
		}
		if _, err := tx.NewInsert().Model(approval).Returning("*").Exec(ctx); err != nil {
			return err
		}

		action := &ent.ApprovalAction{
			RequestID:   approval.ID,
			Action:      ent.ApprovalActionTypeSubmit,
			ActedBy:     data.RequestedBy,
			ActedByRole: data.RequestedByRole,
			Comment:     data.RequestReason,
			Metadata: map[string]any{
				"source":        "membermanagements.register",
				"management_id": management.ID.String(),
			},
			CreatedAt: now,
		}
		if _, err := tx.NewInsert().Model(action).Returning("*").Exec(ctx); err != nil {
			return err
		}

		result.Member = member
		result.Management = management
		result.Approval = approval
		result.ApprovalAction = action
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetMemberManagementByID(ctx context.Context, id uuid.UUID) (*ent.MemberManagement, error) {
	row := new(ent.MemberManagement)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListMemberManagements(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID) ([]*ent.MemberManagement, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.MemberManagement, 0)
	query := s.db.NewSelect().Model(&items)

	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}
	if memberID != nil {
		query.Where("member_id = ?", *memberID)
	}
	if departmentID != nil {
		query.Where("department_id = ?", *departmentID)
	}

	if err := req.SetSearchBy(query, []string{"employee_code", "position"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "employee_code", "start_work_date", "position", "is_active"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateMemberManagementByID(ctx context.Context, id uuid.UUID, data *ent.MemberManagementUpdate) (*ent.MemberManagement, error) {
	query := s.db.NewUpdate().
		Model(&ent.MemberManagement{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.MemberID != nil {
		query.Set("member_id = ?", *data.MemberID)
	}
	if data.EmployeeCode != nil {
		query.Set("employee_code = ?", *data.EmployeeCode)
	}
	if data.Position != nil {
		query.Set("position = ?", *data.Position)
	}
	if data.StartWorkDate != nil {
		query.Set("start_work_date = ?", *data.StartWorkDate)
	}
	if data.DepartmentID != nil {
		query.Set("department_id = ?", *data.DepartmentID)
	}
	if data.IsActive != nil {
		query.Set("is_active = ?", *data.IsActive)
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

	return s.GetMemberManagementByID(ctx, id)
}

func (s *Service) SoftDeleteMemberManagementByID(ctx context.Context, id uuid.UUID) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model(&ent.MemberManagement{}).
			Set("is_active = ?", false).
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

		_, err = tx.NewDelete().Model(&ent.MemberManagement{}).Where("id = ?", id).Exec(ctx)
		return err
	})
}

package entities

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ entitiesinf.MemberManagementEntity = (*Service)(nil)

func isMissingSchoolDepartmentIDColumnError(err error) bool {
	if err == nil {
		return false
	}
	errText := err.Error()
	return strings.Contains(errText, "school_department_id") && strings.Contains(strings.ToLower(errText), "does not exist")
}

func (s *Service) CreateMemberManagement(ctx context.Context, data *ent.MemberManagement) (*ent.MemberManagement, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) RegisterManagement(ctx context.Context, data *ent.ManagementRegistrationInput) (*ent.ManagementRegistrationResult, error) {
	result := new(ent.ManagementRegistrationResult)

	hasSchoolDepartmentColumn, err := s.hasMemberManagementSchoolDepartmentColumn(ctx)
	if err != nil {
		return nil, err
	}

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		member := &ent.Member{
			ID:        uuid.New(),
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

		var management *ent.MemberManagement
		if hasSchoolDepartmentColumn {
			management = &ent.MemberManagement{
				ID:                 uuid.New(),
				MemberID:           member.ID,
				EmployeeCode:       data.ManagementEmployeeCode,
				GenderID:           data.ManagementGenderID,
				PrefixID:           data.ManagementPrefixID,
				FirstName:          data.ManagementFirstName,
				LastName:           data.ManagementLastName,
				Phone:              data.ManagementPhone,
				Position:           data.ManagementPosition,
				StartWorkDate:      data.ManagementStartWorkDate,
				SchoolDepartmentID: data.ManagementSchoolDepartmentID,
				DepartmentID:       data.ManagementDepartmentID,
				IsActive:           data.ManagementIsActive,
			}
			if _, err := tx.NewInsert().Model(management).Returning("*").Exec(ctx); err != nil {
				return err
			}
		} else {
			legacyManagement := &memberManagementLegacyRow{
				ID:            uuid.New(),
				MemberID:      member.ID,
				EmployeeCode:  data.ManagementEmployeeCode,
				GenderID:      data.ManagementGenderID,
				PrefixID:      data.ManagementPrefixID,
				FirstName:     data.ManagementFirstName,
				LastName:      data.ManagementLastName,
				Phone:         data.ManagementPhone,
				Position:      data.ManagementPosition,
				StartWorkDate: data.ManagementStartWorkDate,
				DepartmentID:  data.ManagementDepartmentID,
				IsActive:      data.ManagementIsActive,
			}
			if _, err := tx.NewInsert().Model(legacyManagement).Returning("*").Exec(ctx); err != nil {
				return err
			}

			management = &ent.MemberManagement{
				ID:                 legacyManagement.ID,
				MemberID:           legacyManagement.MemberID,
				EmployeeCode:       legacyManagement.EmployeeCode,
				GenderID:           legacyManagement.GenderID,
				PrefixID:           legacyManagement.PrefixID,
				FirstName:          legacyManagement.FirstName,
				LastName:           legacyManagement.LastName,
				Phone:              legacyManagement.Phone,
				Position:           legacyManagement.Position,
				StartWorkDate:      legacyManagement.StartWorkDate,
				SchoolDepartmentID: data.ManagementSchoolDepartmentID,
				DepartmentID:       legacyManagement.DepartmentID,
				IsActive:           legacyManagement.IsActive,
				CreatedAt:          legacyManagement.CreatedAt,
				UpdatedAt:          legacyManagement.UpdatedAt,
				DeletedAt:          legacyManagement.DeletedAt,
			}
		}

		now := time.Now()
		subjectID := management.ID
		approval := &ent.ApprovalRequest{
			ID:              uuid.New(),
			RequestType:     "management_registration",
			SubjectType:     "member_management",
			SubjectID:       &subjectID,
			RequestedBy:     data.RequestedBy,
			RequestedByRole: data.RequestedByRole,
			Payload: map[string]any{
				"member_id":            member.ID.String(),
				"management_id":        management.ID.String(),
				"employee_code":        management.EmployeeCode,
				"school_department_id": management.SchoolDepartmentID.String(),
				"department_id":        management.DepartmentID.String(),
				"target_role":          string(member.Role),
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
			ID:          uuid.New(),
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
	hasSchoolDepartmentColumn, err := s.hasMemberManagementSchoolDepartmentColumn(ctx)
	if err != nil {
		return nil, err
	}

	if !hasSchoolDepartmentColumn {
		legacyRow := new(memberManagementLegacyRow)
		if err := s.db.NewSelect().
			Model(legacyRow).
			ColumnExpr("mm.*").
			ColumnExpr("sd.id AS school_department_id").
			Join("JOIN members AS m ON m.id = mm.member_id").
			Join("LEFT JOIN school_departments AS sd ON sd.school_id = m.school_id AND sd.department_id = mm.department_id AND sd.deleted_at IS NULL").
			Where("mm.id = ?", id).
			Where("m.role = ?", ent.MemberRoleStaff).
			Where("m.deleted_at IS NULL").
			Scan(ctx); err != nil {
			return nil, err
		}

		schoolDepartmentID := legacyRow.DepartmentID
		if legacyRow.SchoolDepartmentID != nil {
			schoolDepartmentID = *legacyRow.SchoolDepartmentID
		}

		return &ent.MemberManagement{
			ID:                 legacyRow.ID,
			MemberID:           legacyRow.MemberID,
			EmployeeCode:       legacyRow.EmployeeCode,
			GenderID:           legacyRow.GenderID,
			PrefixID:           legacyRow.PrefixID,
			FirstName:          legacyRow.FirstName,
			LastName:           legacyRow.LastName,
			Phone:              legacyRow.Phone,
			Position:           legacyRow.Position,
			StartWorkDate:      legacyRow.StartWorkDate,
			SchoolDepartmentID: schoolDepartmentID,
			DepartmentID:       legacyRow.DepartmentID,
			IsActive:           legacyRow.IsActive,
			CreatedAt:          legacyRow.CreatedAt,
			UpdatedAt:          legacyRow.UpdatedAt,
			DeletedAt:          legacyRow.DeletedAt,
		}, nil
	}

	row := new(ent.MemberManagement)
	if err := s.db.NewSelect().
		Model(row).
		Join("JOIN members AS m ON m.id = mm.member_id").
		Where("mm.id = ?", id).
		Where("m.role = ?", ent.MemberRoleStaff).
		Where("m.deleted_at IS NULL").
		Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

type memberManagementLegacyRow struct {
	bun.BaseModel `bun:"table:member_managements,alias:mm"`

	ID                 uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	MemberID           uuid.UUID  `bun:"member_id,notnull,type:uuid"`
	EmployeeCode       string     `bun:"employee_code,notnull,type:varchar(50)"`
	GenderID           *uuid.UUID `bun:"gender_id,scanonly"`
	PrefixID           *uuid.UUID `bun:"prefix_id,scanonly"`
	FirstName          *string    `bun:"first_name,scanonly"`
	LastName           *string    `bun:"last_name,scanonly"`
	Phone              *string    `bun:"phone,scanonly"`
	Position           string     `bun:"position,notnull,type:varchar(255)"`
	StartWorkDate      time.Time  `bun:"start_work_date,notnull,type:date"`
	DepartmentID       uuid.UUID  `bun:"department_id,notnull,type:uuid"`
	SchoolDepartmentID *uuid.UUID `bun:"school_department_id,scanonly"`
	IsActive           bool       `bun:"is_active,notnull,default:true"`
	CreatedAt          time.Time  `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time  `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt          *time.Time `bun:"deleted_at,soft_delete"`
}

func (s *Service) hasMemberManagementSchoolDepartmentColumn(ctx context.Context) (bool, error) {
	var probe bool
	err := s.db.NewSelect().
		TableExpr("member_managements AS mm").
		ColumnExpr("mm.school_department_id IS NOT NULL").
		Limit(1).
		Scan(ctx, &probe)
	if err == nil || err == sql.ErrNoRows {
		return true, nil
	}

	if isMissingSchoolDepartmentIDColumnError(err) {
		return false, nil
	}

	return false, err
}

func (s *Service) ListMemberManagements(ctx context.Context, req *base.RequestPaginate, isActive *bool, memberID *uuid.UUID, departmentID *uuid.UUID, schoolDepartmentID *uuid.UUID) ([]*ent.MemberManagement, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	hasSchoolDepartmentColumn, err := s.hasMemberManagementSchoolDepartmentColumn(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !hasSchoolDepartmentColumn {
		if schoolDepartmentID != nil {
			schoolDepartment, err := s.GetSchoolDepartmentByID(ctx, *schoolDepartmentID)
			if err != nil {
				return nil, nil, err
			}
			if departmentID != nil && *departmentID != schoolDepartment.DepartmentID {
				return nil, nil, fmt.Errorf("school-department-filter-mismatch")
			}
			departmentID = &schoolDepartment.DepartmentID
		}

		legacyRows := make([]*memberManagementLegacyRow, 0)
		legacyQuery := s.db.NewSelect().
			Model(&legacyRows).
			ColumnExpr("mm.*").
			ColumnExpr("sd.id AS school_department_id").
			Join("JOIN members AS m ON m.id = mm.member_id").
			Join("LEFT JOIN school_departments AS sd ON sd.school_id = m.school_id AND sd.department_id = mm.department_id AND sd.deleted_at IS NULL").
			Where("m.role = ?", ent.MemberRoleStaff).
			Where("m.deleted_at IS NULL")

		if isActive != nil {
			legacyQuery.Where("mm.is_active = ?", *isActive)
		}
		if memberID != nil {
			legacyQuery.Where("mm.member_id = ?", *memberID)
		}
		if departmentID != nil {
			legacyQuery.Where("mm.department_id = ?", *departmentID)
		}

		if err := req.SetSearchBy(legacyQuery, []string{"employee_code", "position"}); err != nil {
			return nil, nil, err
		}

		if err := applyMemberManagementSortOrder(req, legacyQuery); err != nil {
			return nil, nil, err
		}

		req.SetOffsetLimit(legacyQuery)
		total, err := legacyQuery.ScanAndCount(ctx)
		if err != nil {
			return nil, nil, err
		}

		items := make([]*ent.MemberManagement, 0, len(legacyRows))
		for _, row := range legacyRows {
			if row == nil {
				continue
			}
			schoolDepartmentID := row.DepartmentID
			if row.SchoolDepartmentID != nil {
				schoolDepartmentID = *row.SchoolDepartmentID
			}
			items = append(items, &ent.MemberManagement{
				ID:                 row.ID,
				MemberID:           row.MemberID,
				EmployeeCode:       row.EmployeeCode,
				GenderID:           row.GenderID,
				PrefixID:           row.PrefixID,
				FirstName:          row.FirstName,
				LastName:           row.LastName,
				Phone:              row.Phone,
				Position:           row.Position,
				StartWorkDate:      row.StartWorkDate,
				SchoolDepartmentID: schoolDepartmentID,
				DepartmentID:       row.DepartmentID,
				IsActive:           row.IsActive,
				CreatedAt:          row.CreatedAt,
				UpdatedAt:          row.UpdatedAt,
				DeletedAt:          row.DeletedAt,
			})
		}

		return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
	}

	items := make([]*ent.MemberManagement, 0)
	query := s.db.NewSelect().
		Model(&items).
		Join("JOIN members AS m ON m.id = mm.member_id").
		Where("m.role = ?", ent.MemberRoleStaff).
		Where("m.deleted_at IS NULL")

	if isActive != nil {
		query.Where("mm.is_active = ?", *isActive)
	}
	if memberID != nil {
		query.Where("mm.member_id = ?", *memberID)
	}
	if departmentID != nil {
		query.Where("mm.department_id = ?", *departmentID)
	}
	if schoolDepartmentID != nil {
		query.Where("mm.school_department_id = ?", *schoolDepartmentID)
	}

	if err := req.SetSearchBy(query, []string{"employee_code", "position"}); err != nil {
		return nil, nil, err
	}

	if err := applyMemberManagementSortOrder(req, query); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func applyMemberManagementSortOrder(req *base.RequestPaginate, query *bun.SelectQuery) error {
	if req == nil {
		query.OrderExpr("mm.created_at DESC")
		return nil
	}

	order := strings.ToUpper(strings.TrimSpace(req.OrderBy))
	if order != "DESC" {
		order = "ASC"
	}

	sortKey := strings.TrimSpace(req.SortBy)
	if sortKey == "" {
		query.OrderExpr("mm.created_at DESC")
		return nil
	}

	columnBySortKey := map[string]string{
		"created_at":      "mm.created_at",
		"employee_code":   "mm.employee_code",
		"start_work_date": "mm.start_work_date",
		"position":        "mm.position",
		"is_active":       "mm.is_active",
	}

	column, ok := columnBySortKey[sortKey]
	if !ok {
		return base.ErrInvalidSort
	}

	query.OrderExpr(column + " " + order)
	return nil
}

func (s *Service) UpdateMemberManagementByID(ctx context.Context, id uuid.UUID, data *ent.MemberManagementUpdate) (*ent.MemberManagement, error) {
	hasSchoolDepartmentColumn, err := s.hasMemberManagementSchoolDepartmentColumn(ctx)
	if err != nil {
		return nil, err
	}

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
	if data.GenderID != nil {
		query.Set("gender_id = ?", *data.GenderID)
	}
	if data.PrefixID != nil {
		query.Set("prefix_id = ?", *data.PrefixID)
	}
	if data.FirstName != nil {
		query.Set("first_name = ?", *data.FirstName)
	}
	if data.LastName != nil {
		query.Set("last_name = ?", *data.LastName)
	}
	if data.Phone != nil {
		query.Set("phone = ?", *data.Phone)
	}
	if data.Position != nil {
		query.Set("position = ?", *data.Position)
	}
	if data.StartWorkDate != nil {
		query.Set("start_work_date = ?", *data.StartWorkDate)
	}
	if hasSchoolDepartmentColumn && data.SchoolDepartmentID != nil {
		query.Set("school_department_id = ?", *data.SchoolDepartmentID)
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

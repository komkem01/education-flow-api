package entities

import (
	"context"
	"database/sql"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentEnrollmentEntity = (*Service)(nil)

func (s *Service) CreateStudentEnrollment(ctx context.Context, data *ent.StudentEnrollment) (*ent.StudentEnrollment, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetStudentEnrollmentByID(ctx context.Context, id uuid.UUID) (*ent.StudentEnrollment, error) {
	row := new(ent.StudentEnrollment)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListStudentEnrollments(ctx context.Context, req *base.RequestPaginate, studentID *uuid.UUID, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, status *ent.StudentEnrollmentStatus, enrollmentType *ent.EnrollmentType) ([]*ent.StudentEnrollment, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.StudentEnrollment, 0)
	query := s.db.NewSelect().Model(&items)

	if studentID != nil {
		query.Where("student_id = ?", *studentID)
	}
	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if academicYearID != nil {
		query.Where("academic_year_id = ?", *academicYearID)
	}
	if classroomID != nil {
		query.Where("classroom_id = ?", *classroomID)
	}
	if status != nil {
		query.Where("status = ?", *status)
	}
	if enrollmentType != nil {
		query.Where("enrollment_type = ?", *enrollmentType)
	}

	if err := req.SetSearchBy(query, []string{"status", "enrollment_type", "roll_no", "exit_note", "approval_note"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("created_at DESC")
	}
	if err := req.SetSortOrder(query, []string{"created_at", "enrolled_at", "exited_at", "status", "enrollment_type", "roll_no"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateStudentEnrollmentByID(ctx context.Context, id uuid.UUID, data *ent.StudentEnrollmentUpdate) (*ent.StudentEnrollment, error) {
	query := s.db.NewUpdate().
		Model(&ent.StudentEnrollment{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.StudentID != nil {
		query.Set("student_id = ?", *data.StudentID)
	}
	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.AcademicYearID != nil {
		query.Set("academic_year_id = ?", *data.AcademicYearID)
	}
	if data.ClassroomID != nil {
		query.Set("classroom_id = ?", *data.ClassroomID)
	}
	if data.EnrolledAt != nil {
		query.Set("enrolled_at = ?", *data.EnrolledAt)
	}
	if data.ExitedAt != nil {
		query.Set("exited_at = ?", *data.ExitedAt)
	}
	if data.Status != nil {
		query.Set("status = ?", *data.Status)
	}
	if data.EnrollmentType != nil {
		query.Set("enrollment_type = ?", *data.EnrollmentType)
	}
	if data.ExitReason != nil {
		query.Set("exit_reason = ?", *data.ExitReason)
	}
	if data.ExitNote != nil {
		query.Set("exit_note = ?", *data.ExitNote)
	}
	if data.PreviousEnrollmentID != nil {
		query.Set("previous_enrollment_id = ?", *data.PreviousEnrollmentID)
	}
	if data.RollNo != nil {
		query.Set("roll_no = ?", *data.RollNo)
	}
	if data.ApprovedBy != nil {
		query.Set("approved_by = ?", *data.ApprovedBy)
	}
	if data.ApprovedAt != nil {
		query.Set("approved_at = ?", *data.ApprovedAt)
	}
	if data.ApprovalNote != nil {
		query.Set("approval_note = ?", *data.ApprovalNote)
	}
	if data.CreatedBy != nil {
		query.Set("created_by = ?", *data.CreatedBy)
	}
	if data.UpdatedBy != nil {
		query.Set("updated_by = ?", *data.UpdatedBy)
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

	return s.GetStudentEnrollmentByID(ctx, id)
}

func (s *Service) SoftDeleteStudentEnrollmentByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.StudentEnrollment{}).Where("id = ?", id).Exec(ctx)
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

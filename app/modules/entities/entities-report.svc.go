package entities

import (
	"context"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.ReportEntity = (*Service)(nil)

func (s *Service) CountMembersByRole(ctx context.Context, schoolID uuid.UUID, role ent.MemberRole, isActive *bool) (int64, error) {
	query := s.db.NewSelect().Model((*ent.Member)(nil)).ColumnExpr("COUNT(*)")
	query.Where("school_id = ?", schoolID)
	query.Where("role = ?", role)
	if isActive != nil {
		query.Where("is_active = ?", *isActive)
	}

	var total int64
	if err := query.Scan(ctx, &total); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Service) CountSubjectsBySchool(ctx context.Context, schoolID uuid.UUID) (int64, error) {
	query := s.db.NewSelect().Model((*ent.Subject)(nil)).ColumnExpr("COUNT(*)")
	query.Where("school_id = ?", schoolID)

	var total int64
	if err := query.Scan(ctx, &total); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Service) CountAttendanceByStatus(ctx context.Context, schoolID uuid.UUID, academicYearID *uuid.UUID, status *ent.AttendanceStatus) (int64, error) {
	query := s.db.NewSelect().
		TableExpr("attendance_records AS ar").
		ColumnExpr("COUNT(*)").
		Join("JOIN attendance_sessions AS ats ON ats.id = ar.session_id").
		Where("ar.deleted_at IS NULL").
		Where("ats.deleted_at IS NULL").
		Where("ats.school_id = ?", schoolID)

	if academicYearID != nil {
		query.Where("ats.academic_year_id = ?", *academicYearID)
	}

	if status != nil {
		query.Where("ar.status = ?", *status)
	}

	var total int64
	if err := query.Scan(ctx, &total); err != nil {
		return 0, err
	}

	return total, nil
}

func (s *Service) ListAcademicYearsBySchool(ctx context.Context, schoolID uuid.UUID) ([]*ent.AcademicYear, error) {
	items := make([]*ent.AcademicYear, 0)
	if err := s.db.NewSelect().
		Model(&items).
		Where("school_id = ?", schoolID).
		Order("start_date DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

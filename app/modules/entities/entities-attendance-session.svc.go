package entities

import (
	"context"
	"database/sql"
	"time"

	"eduflow/app/modules/entities/ent"
	entitiesinf "eduflow/app/modules/entities/inf"
	"eduflow/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.AttendanceSessionEntity = (*Service)(nil)

func (s *Service) CreateAttendanceSession(ctx context.Context, data *ent.AttendanceSession) (*ent.AttendanceSession, error) {
	if _, err := s.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) GetAttendanceSessionByID(ctx context.Context, id uuid.UUID) (*ent.AttendanceSession, error) {
	row := new(ent.AttendanceSession)
	if err := s.db.NewSelect().Model(row).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListAttendanceSessions(ctx context.Context, req *base.RequestPaginate, schoolID *uuid.UUID, academicYearID *uuid.UUID, classroomID *uuid.UUID, subjectID *uuid.UUID, teacherID *uuid.UUID, mode *ent.AttendanceMode, sessionDateFrom *string, sessionDateTo *string) ([]*ent.AttendanceSession, *base.ResponsePaginate, error) {
	if req == nil {
		req = &base.RequestPaginate{}
	}

	items := make([]*ent.AttendanceSession, 0)
	query := s.db.NewSelect().Model(&items)

	if schoolID != nil {
		query.Where("school_id = ?", *schoolID)
	}
	if academicYearID != nil {
		query.Where("academic_year_id = ?", *academicYearID)
	}
	if classroomID != nil {
		query.Where("classroom_id = ?", *classroomID)
	}
	if subjectID != nil {
		query.Where("subject_id = ?", *subjectID)
	}
	if teacherID != nil {
		query.Where("teacher_id = ?", *teacherID)
	}
	if mode != nil {
		query.Where("mode = ?", *mode)
	}
	if sessionDateFrom != nil && *sessionDateFrom != "" {
		if parsed, err := time.Parse("2006-01-02", *sessionDateFrom); err == nil {
			query.Where("session_date >= ?", parsed)
		}
	}
	if sessionDateTo != nil && *sessionDateTo != "" {
		if parsed, err := time.Parse("2006-01-02", *sessionDateTo); err == nil {
			query.Where("session_date <= ?", parsed)
		}
	}

	if err := req.SetSearchBy(query, []string{"mode", "note"}); err != nil {
		return nil, nil, err
	}

	if req.SortBy == "" {
		query.Order("session_date DESC").Order("period_no DESC")
	}
	if err := req.SetSortOrder(query, []string{"session_date", "period_no", "mode", "created_at"}); err != nil {
		return nil, nil, err
	}

	req.SetOffsetLimit(query)
	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, nil, err
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateAttendanceSessionByID(ctx context.Context, id uuid.UUID, data *ent.AttendanceSessionUpdate) (*ent.AttendanceSession, error) {
	query := s.db.NewUpdate().
		Model(&ent.AttendanceSession{}).
		Where("id = ?", id).
		Set("updated_at = now()")

	if data.SchoolID != nil {
		query.Set("school_id = ?", *data.SchoolID)
	}
	if data.AcademicYearID != nil {
		query.Set("academic_year_id = ?", *data.AcademicYearID)
	}
	if data.ClassroomID != nil {
		query.Set("classroom_id = ?", *data.ClassroomID)
	}
	if data.SubjectID != nil {
		query.Set("subject_id = ?", *data.SubjectID)
	}
	if data.TeacherID != nil {
		query.Set("teacher_id = ?", *data.TeacherID)
	}
	if data.SessionDate != nil {
		query.Set("session_date = ?", *data.SessionDate)
	}
	if data.PeriodNo != nil {
		query.Set("period_no = ?", *data.PeriodNo)
	}
	if data.Mode != nil {
		query.Set("mode = ?", *data.Mode)
	}
	if data.StartedAt != nil {
		query.Set("started_at = ?", *data.StartedAt)
	}
	if data.ClosedAt != nil {
		query.Set("closed_at = ?", *data.ClosedAt)
	}
	if data.Note != nil {
		query.Set("note = ?", *data.Note)
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

	return s.GetAttendanceSessionByID(ctx, id)
}

func (s *Service) SoftDeleteAttendanceSessionByID(ctx context.Context, id uuid.UUID) error {
	res, err := s.db.NewDelete().Model(&ent.AttendanceSession{}).Where("id = ?", id).Exec(ctx)
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

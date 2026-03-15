package attendancerecords

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

type BulkMarkResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
}

func (s *Service) BulkMark(ctx context.Context, req *BulkMarkRequest) (*BulkMarkResult, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "attendancerecords.service.bulk_mark")
	defer span.End()

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	sourceVal, ok := parseAttendanceSource(req.Source)
	if !ok {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	markedByVal, err := parseOptionalUUID(req.MarkedBy)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
	}

	result := &BulkMarkResult{}

	for _, item := range req.Items {
		enrollmentID, err := uuid.Parse(item.EnrollmentID)
		if err != nil {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}

		statusVal, ok := parseAttendanceStatus(item.Status)
		if !ok {
			return nil, fmt.Errorf("%w", ErrAttendanceRecordConditionFail)
		}

		existing, err := s.db.GetAttendanceRecordBySessionAndEnrollment(ctx, sessionID, enrollmentID)
		if err != nil {
			if !isNoRows(err) {
				return nil, normalizeServiceError(err)
			}

			created, err := s.db.CreateAttendanceRecord(ctx, &ent.AttendanceRecord{
				SessionID:    sessionID,
				EnrollmentID: enrollmentID,
				Status:       statusVal,
				Source:       sourceVal,
				MarkedAt:     time.Now(),
				Remark:       item.Remark,
				MarkedBy:     markedByVal,
			})
			if err != nil {
				return nil, normalizeServiceError(err)
			}

			_, _ = s.logDB.CreateAttendanceRecordLog(ctx, &ent.AttendanceRecordLog{
				RecordID:  created.ID,
				OldStatus: nil,
				NewStatus: statusVal,
				ChangedBy: markedByVal,
				ChangedAt: time.Now(),
				Reason:    req.Reason,
			})

			result.Created++
			continue
		}

		if existing.Status == statusVal && existing.Source == sourceVal {
			continue
		}

		old := existing.Status
		_, err = s.db.UpdateAttendanceRecordByID(ctx, existing.ID, &ent.AttendanceRecordUpdate{
			Status:   &statusVal,
			Source:   &sourceVal,
			MarkedAt: ptrTime(time.Now()),
			Remark:   item.Remark,
			MarkedBy: markedByVal,
		})
		if err != nil {
			return nil, normalizeServiceError(err)
		}

		_, _ = s.logDB.CreateAttendanceRecordLog(ctx, &ent.AttendanceRecordLog{
			RecordID:  existing.ID,
			OldStatus: &old,
			NewStatus: statusVal,
			ChangedBy: markedByVal,
			ChangedAt: time.Now(),
			Reason:    req.Reason,
		})

		result.Updated++
	}

	return result, nil
}

func isNoRows(err error) bool {
	return err != nil && (err == sql.ErrNoRows || fmt.Sprintf("%v", err) == sql.ErrNoRows.Error())
}

func ptrTime(v time.Time) *time.Time {
	return &v
}

package schoolannouncements

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrSchoolAnnouncementNotFound      = errors.New("school-announcement-not-found")
	ErrSchoolAnnouncementInvalidRole   = errors.New("school-announcement-invalid-role")
	ErrSchoolAnnouncementInvalidUpdate = errors.New("school-announcement-empty-update-payload")
	ErrSchoolAnnouncementUnauthorized  = errors.New("school-announcement-unauthorized")
	ErrSchoolAnnouncementConditionFail = errors.New("school-announcement-condition-fail")
)

func normalizeServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", ErrSchoolAnnouncementNotFound, err)
	}
	return err
}

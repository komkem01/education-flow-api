package schoolannouncements

import (
	"time"

	"eduflow/app/modules/entities/ent"
)

func isValidStatus(status ent.SchoolAnnouncementStatus) bool {
	switch status {
	case ent.SchoolAnnouncementStatusDraft, ent.SchoolAnnouncementStatusPublished, ent.SchoolAnnouncementStatusExpired:
		return true
	default:
		return false
	}
}

func isValidTargetRole(role string) bool {
	switch role {
	case "admin", "staff", "teacher", "student", "parent", "":
		return true
	default:
		return false
	}
}

func parseDatePtr(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	if *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *raw)
	if err == nil {
		return &parsed, nil
	}
	parsed, err = time.Parse(time.RFC3339, *raw)
	if err == nil {
		v := parsed.UTC()
		return &v, nil
	}
	return nil, err
}

package reports

import (
	"context"
	"fmt"

	"eduflow/app/modules/entities/ent"
	"eduflow/app/utils"

	"github.com/google/uuid"
)

type AcademicYearFilterItem struct {
	ID    uuid.UUID `json:"id"`
	Year  string    `json:"year"`
	Term  string    `json:"term"`
	Label string    `json:"label"`
}

type FiltersResponse struct {
	AcademicYears []*AcademicYearFilterItem `json:"academic_years"`
	Semesters     []int                     `json:"semesters"`
}

func (s *Service) Filters(ctx context.Context, schoolID uuid.UUID) (*FiltersResponse, error) {
	ctx, span, _ := utils.NewLogSpan(ctx, s.tracer, "reports.service.filters")
	defer span.End()

	years, err := s.db.ListAcademicYearsBySchool(ctx, schoolID)
	if err != nil {
		return nil, fmt.Errorf("list-academic-years: %w", err)
	}

	items := make([]*AcademicYearFilterItem, 0, len(years))
	for _, year := range years {
		if year == nil {
			continue
		}
		items = append(items, &AcademicYearFilterItem{
			ID:    year.ID,
			Year:  year.Year,
			Term:  "",
			Label: year.Year,
		})
	}

	return &FiltersResponse{
		AcademicYears: items,
		Semesters:     []int{1, 2},
	}, nil
}

func fallbackSchoolID(user *ent.Member) uuid.UUID {
	if user == nil {
		return uuid.Nil
	}
	return user.SchoolID
}

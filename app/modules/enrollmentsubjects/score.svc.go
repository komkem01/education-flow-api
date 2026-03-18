package enrollmentsubjects

import (
	"math"

	"eduflow/app/modules/entities/ent"
)

func enrichScoreSummary(item *ent.EnrollmentSubject) {
	if item == nil {
		return
	}

	total := computeTotalScore(item.MidtermScore, item.FinalScore, item.ActivityScore)
	item.TotalScore = total
	item.GradeNumeric = computeGradeNumeric(total)
	item.GradeLetter = computeGradeLetter(total)
}

func enrichScoreSummaryList(items []*ent.EnrollmentSubject) {
	for _, item := range items {
		enrichScoreSummary(item)
	}
}

func computeTotalScore(midterm, final, activity *float64) *float64 {
	if midterm == nil && final == nil && activity == nil {
		return nil
	}

	total := valueOrZero(midterm) + valueOrZero(final) + valueOrZero(activity)
	total = math.Round(total*100) / 100
	return &total
}

func computeGradeNumeric(total *float64) *string {
	if total == nil {
		return nil
	}

	var grade string
	score := *total
	switch {
	case score >= 80:
		grade = "4"
	case score >= 75:
		grade = "3.5"
	case score >= 70:
		grade = "3"
	case score >= 65:
		grade = "2.5"
	case score >= 60:
		grade = "2"
	case score >= 55:
		grade = "1.5"
	case score >= 50:
		grade = "1"
	default:
		grade = "0"
	}

	return &grade
}

func computeGradeLetter(total *float64) *string {
	if total == nil {
		return nil
	}

	var grade string
	score := *total
	switch {
	case score >= 80:
		grade = "A"
	case score >= 75:
		grade = "B+"
	case score >= 70:
		grade = "B"
	case score >= 65:
		grade = "C+"
	case score >= 60:
		grade = "C"
	case score >= 55:
		grade = "D+"
	case score >= 50:
		grade = "D"
	default:
		grade = "F"
	}

	return &grade
}

func valueOrZero(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

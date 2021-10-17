package domain

import (
	"context"
	"errors"
)

var ErrStudyResultNotFound = errors.New("StudyResult not found")

type StudyResultRepository interface {
	FindStudyResults(ctx context.Context, operator Student, workbookID WorkbookID, studyType string) (map[ProblemID]int, error)
	SetResult(ctx context.Context, operator Student, workbookID WorkbookID, studyType string, problemType string, problemID ProblemID, result bool) error
}

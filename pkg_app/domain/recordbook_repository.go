package domain

import (
	"context"
	"errors"
)

var ErrStudyResultNotFound = errors.New("StudyResult not found")

type RecordbookRepository interface {
	FindStudyResults(ctx context.Context, operator Student, workbookID WorkbookID, studyType string) (map[ProblemID]StudyStatus, error)

	SetResult(ctx context.Context, operator Student, workbookID WorkbookID, studyType string, problemType string, problemID ProblemID, result, memorized bool) error
}

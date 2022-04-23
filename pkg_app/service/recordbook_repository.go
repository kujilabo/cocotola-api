//go:generate mockery --output mock --name RecordbookRepository
package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrStudyResultNotFound = errors.New("StudyResult not found")

type RecordbookRepository interface {
	FindStudyResults(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string) (map[domain.ProblemID]domain.StudyStatus, error)

	SetResult(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID, studyType string, problemType string, problemID domain.ProblemID, studyResult, memorized bool) error

	CountMemorizedProblem(ctx context.Context, operator domain.StudentModel, workbookID domain.WorkbookID) (map[string]int, error)
}

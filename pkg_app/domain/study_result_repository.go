package domain

import "context"

type StudyResultRepository interface {
	FindStudyResults(ctx context.Context, operator Student, workbookID WorkbookID, studyType string) (map[ProblemID]int, error)
}

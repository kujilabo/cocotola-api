package domain

import (
	"context"
	"io"
)

type Added int
type Updated int

type ProblemAddProcessor interface {
	AddProblem(ctx context.Context, repo RepositoryFactory, operator Student, workbook Workbook, param ProblemAddParameter) (Added, ProblemID, error)
}

type ProblemUpdateProcessor interface {
	UpdateProblem(ctx context.Context, repo RepositoryFactory, operator Student, workbook Workbook, param ProblemUpdateParameter) (Added, Updated, error)
}

type ProblemRemoveProcessor interface {
	RemoveProblem(ctx context.Context, repo RepositoryFactory, operator Student, problemID ProblemID, version int) error
}

type ProblemImportProcessor interface {
	CreateCSVReader(ctx context.Context, workbookID WorkbookID, reader io.Reader) (ProblemAddParameterIterator, error)
}

type ProblemQuotaProcessor interface {
	// IsExceeded(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	// Increment(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	// Decrement(ctx context.Context, repo RepositoryFactory, operator Student, name string) (bool, error)

	GetUnitForSizeQuota() QuotaUnit

	GetLimitForSizeQuota() int

	GetUnitForUpdateQuota() QuotaUnit

	GetLimitForUpdateQuota() int
}

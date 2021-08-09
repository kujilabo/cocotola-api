package domain

import (
	"context"
	"io"
)

type ProblemAddProcessor interface {
	AddProblem(ctx context.Context, repo RepositoryFactory, operator Student, param ProblemAddParameter) (ProblemID, error)
}

type ProblemRemoveProcessor interface {
	RemoveProblem(ctx context.Context, repo RepositoryFactory, operator Student, problemID ProblemID, version int) error
}

type ProblemImportProcessor interface {
	CreateCSVReader(ctx context.Context, workbookID WorkbookID, problemType string, reader io.Reader) (ProblemAddParameterIterator, error)
}

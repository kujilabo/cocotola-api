package domain

import "context"

type StudyTypeRepository interface {
	FindAllStudyTypes(ctx context.Context) ([]StudyType, error)
}

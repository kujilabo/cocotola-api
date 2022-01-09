package domain

import (
	"context"
	"errors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrAzureTranslationAlreadyExists = errors.New("azure translation already exists")

type AzureTranslation struct {
	Pos        WordPos
	Target     string
	Confidence float64
}

type TranslationSearchCondition struct {
	PageNo   int
	PageSize int
}

type TranslationSearchResult struct {
	TotalCount int64
	Results    [][]AzureTranslation
}

type AzureTranslationRepository interface {
	Add(ctx context.Context, text string, lang app.Lang2, result []AzureTranslation) error

	Find(ctx context.Context, text string, lang app.Lang2) ([]AzureTranslation, error)

	Contain(ctx context.Context, text string, lang app.Lang2) (bool, error)
}

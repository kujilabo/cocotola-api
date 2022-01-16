package domain

import (
	"context"
	"errors"
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrAzureTranslationAlreadyExists = errors.New("azure translation already exists")

type AzureTranslation struct {
	Pos        WordPos
	Target     string
	Confidence float64
}

func (t *AzureTranslation) ToTranslation(lang app.Lang2, text string) (Translation, error) {
	return NewTranslation(0, 1, time.Now(), time.Now(), text, t.Pos, lang, t.Target, "azure")
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
	Add(ctx context.Context, lang app.Lang2, text string, result []AzureTranslation) error

	Find(ctx context.Context, lang app.Lang2, text string) ([]AzureTranslation, error)

	FindByTextAndPos(ctx context.Context, lang app.Lang2, text string, pos WordPos) (Translation, error)

	FindByFirstLetter(ctx context.Context, lang app.Lang2, firstLetter string) ([]Translation, error)

	Contain(ctx context.Context, lang app.Lang2, text string) (bool, error)
}

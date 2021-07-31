package domain

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type TranslationResult struct {
	Pos        WordPos
	Target     string
	Confidence float64
}

type Translator interface {
	DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]TranslationResult, error)
}

type TranslationSearchCondition struct {
	PageNo   int
	PageSize int
}

type TranslationSearchResult struct {
	TotalCount int64
	Results    [][]TranslationResult
}

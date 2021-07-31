package domain

import (
	"context"
	"errors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrAzureTranslationNotFound = errors.New("azure translation not found")
var ErrAzureTranslationAlreadyExists = errors.New("azure translation already exists")

type AzureTranslationRepository interface {
	Add(ctx context.Context, text string, lang app.Lang2, result []TranslationResult) error

	Find(ctx context.Context, text string, lang app.Lang2) ([]TranslationResult, error)

	Contain(ctx context.Context, text string, lang app.Lang2) (bool, error)
}

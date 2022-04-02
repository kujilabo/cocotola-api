package service

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

type AzureTranslationClient interface {
	DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]AzureTranslation, error)
}

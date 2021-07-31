package gateway

import (
	"context"

	"golang.org/x/xerrors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type azureCachedTranslatorClient struct {
	translator                 domain.Translator
	azureTranslationRepository domain.AzureTranslationRepository
}

func NewAzureCachedTranslatorClient(translator domain.Translator, azureTranslationRepository domain.AzureTranslationRepository) domain.Translator {
	return &azureCachedTranslatorClient{
		translator:                 translator,
		azureTranslationRepository: azureTranslationRepository,
	}

}

func (c *azureCachedTranslatorClient) DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]domain.TranslationResult, error) {
	contained, err := c.azureTranslationRepository.Contain(ctx, text, toLang)
	if err != nil {
		return nil, err
	}
	if contained {
		return c.azureTranslationRepository.Find(ctx, text, toLang)
	}

	result, err := c.translator.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	if err := c.azureTranslationRepository.Add(ctx, text, toLang, result); err != nil {
		return nil, xerrors.Errorf("failed to add auzre_translation. err: %w", err)
	}

	return result, nil
}

package domain

import (
	"context"
	"errors"
	"time"

	"golang.org/x/xerrors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
)

var ErrTranslationNotFound = errors.New("translation not found")

type Translator interface {
	DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]Translation, error)
	DictionaryLookupWithPos(ctx context.Context, text string, fromLang, toLang app.Lang2, pos WordPos) (Translation, error)
	FindWords(ctx context.Context, letter string) ([]Translation, error)
}

type translator struct {
	rf          RepositoryFactory
	azureClient AzureTranslationClient
}

func NewTranslatior(rf RepositoryFactory, azureClient AzureTranslationClient) (Translator, error) {
	return &translator{
		rf:          rf,
		azureClient: azureClient,
	}, nil
}

// func (t *translator) selectMaxConfidenceTranslation(ctx context.Context, in []AzureTranslation, pos WordPos) (AzureTranslation, error) {
// 	found := false
// 	var result AzureTranslation
// 	for _, i := range in {
// 		if i.Pos == pos && i.Confidence > result.Confidence {
// 			found = true
// 			result = i
// 		}
// 	}
// 	if !found {
// 		return result, ErrTranslationNotFound
// 	}
// 	return result, nil
// }

func (t *translator) selectMaxConfidenceTranslations(ctx context.Context, in []AzureTranslation) (map[WordPos]AzureTranslation, error) {
	results := make(map[WordPos]AzureTranslation)
	for _, i := range in {
		if _, ok := results[i.Pos]; !ok {
			results[i.Pos] = i
		} else if i.Confidence > results[i.Pos].Confidence {
			results[i.Pos] = i
		}
	}
	return results, nil
}

func (t *translator) customDictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]Translation, error) {
	// repo, err := t.rf.NewAzureTranslationRepository()
	// if err != nil {
	// 	return nil, err
	// }
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customContained, err := customRepo.Contain(ctx, text, toLang)
	if err != nil {
		return nil, err
	}
	if !customContained {
		return nil, ErrTranslationNotFound
	}

	customResults, err := customRepo.FindByText(ctx, text, toLang)
	if err != nil {
		return nil, err
	}
	return customResults, nil
}

func (t *translator) azureDictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]AzureTranslation, error) {
	// repo, err := t.repo(t.db)
	// if err != nil {
	// 	return nil, err
	// }

	azureRepo, err := t.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureContained, err := azureRepo.Contain(ctx, text, toLang)
	if err != nil {
		return nil, err
	}
	if azureContained {
		azureResults, err := azureRepo.Find(ctx, text, toLang)
		if err != nil {
			return nil, err
		}
		return azureResults, nil
	}

	azureResults, err := t.azureClient.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	if err := azureRepo.Add(ctx, text, toLang, azureResults); err != nil {
		return nil, xerrors.Errorf("failed to add auzre_translation. err: %w", err)
	}

	return azureResults, nil
}

func (t *translator) DictionaryLookup(ctx context.Context, text string, fromLang, toLang app.Lang2) ([]Translation, error) {
	// find translations from custom reopository
	customResults, err := t.customDictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil && !errors.Is(err, ErrTranslationNotFound) {
		return nil, err
	}
	if !errors.Is(err, ErrTranslationNotFound) {
		return customResults, err
	}

	// find translations from azure
	azureResults, err := t.azureDictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}
	azureResultMap, err := t.selectMaxConfidenceTranslations(ctx, azureResults)
	if err != nil {
		return nil, err
	}
	results := make([]Translation, 0)
	for _, v := range azureResultMap {
		result, err := NewTranslation(0, 0, time.Now(), time.Now(), text, v.Pos, fromLang, v.Target)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (t *translator) DictionaryLookupWithPos(ctx context.Context, text string, fromLang, toLang app.Lang2, pos WordPos) (Translation, error) {
	results, err := t.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.GetPos() == pos {
			return r, nil
		}
	}
	return nil, ErrTranslationNotFound
}

func (t *translator) FindWords(ctx context.Context, letter string) ([]Translation, error) {
	return nil, nil
}

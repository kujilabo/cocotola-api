package domain

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
)

var ErrTranslationNotFound = errors.New("translation not found")

type Translator interface {
	DictionaryLookup(ctx context.Context, fromLang, toLang app.Lang2, text string) ([]Translation, error)
	DictionaryLookupWithPos(ctx context.Context, fromLang, toLang app.Lang2, text string, pos WordPos) (Translation, error)
	FindTranslationsByFirstLetter(ctx context.Context, lang app.Lang2, firstLetter string) ([]Translation, error)
	FindTranslationByTextAndPos(ctx context.Context, lang app.Lang2, text string, pos WordPos) (Translation, error)
	FindTranslationByText(ctx context.Context, lang app.Lang2, text string) ([]Translation, error)
	AddTranslation(ctx context.Context, param TranslationAddParameter) error
	UpdateTranslation(ctx context.Context, lang app.Lang2, text string, pos WordPos, param TranslationUpdateParameter) error
	RemoveTranslation(ctx context.Context, lang app.Lang2, text string, pos WordPos) error
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
	customContained, err := customRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if !customContained {
		return nil, ErrTranslationNotFound
	}

	customResults, err := customRepo.FindByText(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	return customResults, nil
}

func (t *translator) azureDictionaryLookup(ctx context.Context, fromLang, toLang app.Lang2, text string) ([]AzureTranslation, error) {
	// repo, err := t.repo(t.db)
	// if err != nil {
	// 	return nil, err
	// }

	azureRepo, err := t.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureContained, err := azureRepo.Contain(ctx, toLang, text)
	if err != nil {
		return nil, err
	}
	if azureContained {
		azureResults, err := azureRepo.Find(ctx, toLang, text)
		if err != nil {
			return nil, err
		}
		return azureResults, nil
	}

	azureResults, err := t.azureClient.DictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil {
		return nil, err
	}

	if err := azureRepo.Add(ctx, toLang, text, azureResults); err != nil {
		return nil, fmt.Errorf("failed to add auzre_translation. err: %w", err)
	}

	return azureResults, nil
}

func (t *translator) DictionaryLookup(ctx context.Context, fromLang, toLang app.Lang2, text string) ([]Translation, error) {
	// find translations from custom reopository
	customResults, err := t.customDictionaryLookup(ctx, text, fromLang, toLang)
	if err != nil && !errors.Is(err, ErrTranslationNotFound) {
		return nil, err
	}
	if !errors.Is(err, ErrTranslationNotFound) {
		return customResults, err
	}

	// find translations from azure
	azureResults, err := t.azureDictionaryLookup(ctx, fromLang, toLang, text)
	if err != nil {
		return nil, err
	}
	azureResultMap, err := t.selectMaxConfidenceTranslations(ctx, azureResults)
	if err != nil {
		return nil, err
	}
	results := make([]Translation, 0)
	for _, v := range azureResultMap {
		result, err := NewTranslation(0, 0, time.Now(), time.Now(), text, v.Pos, fromLang, v.Target, "azure")
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (t *translator) DictionaryLookupWithPos(ctx context.Context, fromLang, toLang app.Lang2, text string, pos WordPos) (Translation, error) {
	results, err := t.DictionaryLookup(ctx, fromLang, toLang, text)
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

func (t *translator) FindTranslationsByFirstLetter(ctx context.Context, lang app.Lang2, firstLetter string) ([]Translation, error) {
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResults, err := customRepo.FindByFirstLetter(ctx, lang, firstLetter)
	if err != nil {
		return nil, err
	}
	azureRepo, err := t.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResults, err := azureRepo.FindByFirstLetter(ctx, lang, firstLetter)
	if err != nil {
		return nil, err
	}

	makeKey := func(text string, pos WordPos) string {
		return text + "_" + strconv.Itoa(int(pos))
	}
	resultMap := make(map[string]Translation)
	for _, c := range customResults {
		key := makeKey(c.GetText(), c.GetPos())
		resultMap[key] = c
	}
	for _, a := range azureResults {
		key := makeKey(a.GetText(), a.GetPos())
		if _, ok := resultMap[key]; !ok {
			resultMap[key] = a
		}
	}

	results := make([]Translation, 0)
	for _, v := range resultMap {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].GetText() < results[j].GetText() })

	return results, nil
}

func (t *translator) FindTranslationByTextAndPos(ctx context.Context, lang app.Lang2, text string, pos WordPos) (Translation, error) {
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResult, err := customRepo.FindByTextAndPos(ctx, lang, text, pos)
	if err == nil {
		return customResult, nil
	}
	if !errors.Is(err, ErrTranslationNotFound) {
		return nil, err
	}

	azureRepo, err := t.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResult, err := azureRepo.FindByTextAndPos(ctx, lang, text, pos)
	if err != nil {
		return nil, err
	}
	return azureResult, nil
}

func (t *translator) FindTranslationByText(ctx context.Context, lang app.Lang2, text string) ([]Translation, error) {
	logger := log.FromContext(ctx)
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	customResults, err := customRepo.FindByText(ctx, lang, text)
	if err != nil {
		return nil, err
	}
	azureRepo, err := t.rf.NewAzureTranslationRepository(ctx)
	if err != nil {
		return nil, err
	}
	azureResults, err := azureRepo.FindByText(ctx, lang, text)
	if err != nil {
		return nil, err
	}

	makeKey := func(text string, pos WordPos) string {
		return text + "_" + strconv.Itoa(int(pos))
	}
	resultMap := make(map[string]Translation)
	for _, c := range customResults {
		key := makeKey(c.GetText(), c.GetPos())
		resultMap[key] = c
	}
	for _, a := range azureResults {
		key := makeKey(a.GetText(), a.GetPos())
		if _, ok := resultMap[key]; !ok {
			resultMap[key] = a
			logger.Infof("translation: %v", a)
		}
	}

	results := make([]Translation, 0)
	for _, v := range resultMap {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].GetPos() < results[j].GetPos() })

	return results, nil
}

func (t *translator) AddTranslation(ctx context.Context, param TranslationAddParameter) error {
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	if _, err := customRepo.Add(ctx, param); err != nil {
		return err
	}
	return nil
}

func (t *translator) UpdateTranslation(ctx context.Context, lang app.Lang2, text string, pos WordPos, param TranslationUpdateParameter) error {
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	translationFound := true
	if _, err := customRepo.FindByTextAndPos(ctx, lang, text, pos); err != nil {
		if errors.Is(err, ErrTranslationNotFound) {
			translationFound = false
		} else {
			return err
		}
	}

	if translationFound {
		if err := customRepo.Update(ctx, lang, text, pos, param); err != nil {
			return err
		}
		return nil
	}

	paramToAdd, err := NewTransalationAddParameter(text, pos, lang, param.GetTranslated())
	if err != nil {
		return err
	}
	if _, err := customRepo.Add(ctx, paramToAdd); err != nil {
		return err
	}
	return nil
}

func (t *translator) RemoveTranslation(ctx context.Context, lang app.Lang2, text string, pos WordPos) error {
	customRepo, err := t.rf.NewCustomTranslationRepository(ctx)
	if err != nil {
		return err
	}

	if err := customRepo.Remove(ctx, lang, text, pos); err != nil {
		return err
	}
	return nil
}

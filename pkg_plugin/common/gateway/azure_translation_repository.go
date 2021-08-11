package gateway

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	libG "github.com/kujilabo/cocotola-api/pkg_lib/gateway"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

type azureTranslationRepository struct {
	db *gorm.DB
}

type azureTranslationEntity struct {
	Text   string
	Lang   string
	Result string
}

func (e *azureTranslationEntity) TableName() string {
	return "azure_translation"
}

func NewAzureTranslationRepository(db *gorm.DB) domain.AzureTranslationRepository {
	return &azureTranslationRepository{
		db: db,
	}
}

func (r *azureTranslationRepository) Add(ctx context.Context, text string, lang app.Lang2, result []domain.TranslationResult) error {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	entity := azureTranslationEntity{
		Text:   text,
		Lang:   lang.String(),
		Result: string(resultBytes),
	}

	if result := r.db.Create(&entity); result.Error != nil {
		return libG.ConvertDuplicatedError(result.Error, domain.ErrAzureTranslationAlreadyExists)
	}

	return nil
}

func (r *azureTranslationRepository) Find(ctx context.Context, text string, lang app.Lang2) ([]domain.TranslationResult, error) {
	entity := azureTranslationEntity{}

	if result := r.db.Where(&azureTranslationEntity{
		Text: text,
		Lang: lang.String(),
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAzureTranslationNotFound
		}

		return nil, result.Error
	}

	result := make([]domain.TranslationResult, 0)
	if err := json.Unmarshal([]byte(entity.Result), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *azureTranslationRepository) FindTranslations(ctx context.Context, param *domain.TranslationSearchCondition) (*domain.TranslationSearchResult, error) {
	limit := param.PageSize
	offset := (param.PageNo - 1) * param.PageSize
	var entities []azureTranslationEntity
	if result := r.db.Limit(limit).Offset(offset).Find(&entities); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	if result := r.db.Model(azureTranslationEntity{}).Count(&count); result.Error != nil {
		return nil, result.Error
	}

	results := make([][]domain.TranslationResult, len(entities))
	for i, e := range entities {
		result := make([]domain.TranslationResult, 0)
		if err := json.Unmarshal([]byte(e.Result), &result); err != nil {
			return nil, err
		}
		results[i] = result
	}

	return &domain.TranslationSearchResult{
		TotalCount: count,
		Results:    results,
	}, nil
}

func (r *azureTranslationRepository) Contain(ctx context.Context, text string, lang app.Lang2) (bool, error) {
	entity := azureTranslationEntity{}

	if result := r.db.Where(&azureTranslationEntity{
		Text: text,
		Lang: lang.String(),
	}).First(&entity); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

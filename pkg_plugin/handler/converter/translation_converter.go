package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/handler/entity"
)

func ToTranslationFindResposne(context context.Context, translations []domain.Translation) (*entity.TranslationFindResponse, error) {

	results := make([]entity.Translation, len(translations))
	for i, t := range translations {
		results[i] = entity.Translation{
			Lang:       t.GetLang().String(),
			Text:       t.GetText(),
			Pos:        int(t.GetPos()),
			Translated: t.GetTranslated(),
			Provider:   t.GetProvider(),
		}
	}

	return &entity.TranslationFindResponse{
		Results: results,
	}, nil
}

func ToTranslationResposne(context context.Context, translation domain.Translation) (*entity.Translation, error) {
	return &entity.Translation{
		Lang:       translation.GetLang().String(),
		Text:       translation.GetText(),
		Pos:        int(translation.GetPos()),
		Translated: translation.GetTranslated(),
		Provider:   translation.GetProvider(),
	}, nil
}

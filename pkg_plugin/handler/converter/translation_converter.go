package converter

import (
	"context"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/handler/entity"
)

func ToTranslationFindResposne(ctx context.Context, translations []domain.Translation) (*entity.TranslationFindResponse, error) {

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

func ToTranslationListResposne(context context.Context, translations []domain.Translation) ([]*entity.Translation, error) {
	results := make([]*entity.Translation, 0)
	for _, t := range translations {
		e := &entity.Translation{
			Lang:       t.GetLang().String(),
			Text:       t.GetText(),
			Pos:        int(t.GetPos()),
			Translated: t.GetTranslated(),
			Provider:   t.GetProvider(),
		}
		results = append(results, e)
	}
	return results, nil
}

func ToTranslationAddParameter(ctx context.Context, param *entity.TranslationAddParameter) (domain.TranslationAddParameter, error) {
	pos, err := domain.NewWordPos(param.Pos)
	if err != nil {
		return nil, err
	}

	lang, err := app.NewLang2(param.Lang)
	if err != nil {
		return nil, err
	}
	return domain.NewTransalationAddParameter(param.Text, pos, lang, param.Translated)
}

func ToTranslationUpdateParameter(ctx context.Context, param *entity.TranslationUpdateParameter) (domain.TranslationUpdateParameter, error) {
	return domain.NewTransaltionUpdateParameter(param.Translated)
}
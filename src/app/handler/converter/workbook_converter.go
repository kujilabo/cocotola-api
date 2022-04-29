package converter

import (
	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/handler/entity"
	"github.com/kujilabo/cocotola-api/src/app/service"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
)

func ToWorkbookSearchResponse(result service.WorkbookSearchResult) (*entity.WorkbookSearchResponse, error) {
	workbooks := make([]*entity.Workbook, len(result.GetResults()))
	for i, w := range result.GetResults() {
		model, err := entity.NewModel(w)
		if err != nil {
			return nil, err
		}

		workbooks[i] = &entity.Workbook{
			Model:        model,
			Name:         w.GetName(),
			Lang2:        w.GetLang2().String(),
			ProblemType:  w.GetProblemType(),
			QuestionText: w.GetQuestionText(),
		}
	}

	e := &entity.WorkbookSearchResponse{
		TotalCount: result.GetTotalCount(),
		Results:    workbooks,
	}
	return e, libD.Validator.Struct(e)
}

func ToWorkbookAddParameter(param *entity.WorkbookAddParameter) (service.WorkbookAddParameter, error) {
	return service.NewWorkbookAddParameter(param.ProblemType, param.Name, domain.Lang2JA, param.QuestionText, map[string]string{
		"audioEnabled": "true",
	})
}

func ToWorkbookUpdateParameter(param *entity.WorkbookUpdateParameter) (service.WorkbookUpdateParameter, error) {
	return service.NewWorkbookUpdateParameter(param.Name, param.QuestionText)
}

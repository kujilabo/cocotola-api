package converter

import (
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
)

func ToWorkbookSearchResponse(result *domain.WorkbookSearchResult) (*entity.WorkbookSearchResponse, error) {
	workbooks := make([]entity.Workbook, len(result.Results))
	for i, w := range result.Results {
		model, err := entity.NewModel(w)
		if err != nil {
			return nil, err
		}

		workbooks[i] = entity.Workbook{
			Model:        model,
			Name:         w.GetName(),
			ProblemType:  w.GetProblemType(),
			QuestionText: w.GetQuestionText(),
		}
	}

	return &entity.WorkbookSearchResponse{
		TotalCount: result.TotalCount,
		Results:    workbooks,
	}, nil
}

func ToWorkbookAddParameter(param *entity.WorkbookAddParameter) (domain.WorkbookAddParameter, error) {
	return domain.NewWorkbookAddParameter(param.ProblemType, param.Name, param.QuestionText)
}

func ToWorkbookUpdateParameter(param *entity.WorkbookUpdateParameter) (domain.WorkbookUpdateParameter, error) {
	return domain.NewWorkbookUpdateParameter(param.Name, param.QuestionText)
}

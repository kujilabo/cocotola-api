package converter

import (
	"context"
	"encoding/json"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

func ToProblemSearchCondition(ctx context.Context, param *entity.ProblemFindParameter, workbookID domain.WorkbookID) (service.ProblemSearchCondition, error) {
	return service.NewProblemSearchCondition(workbookID, param.PageNo, param.PageSize, param.Keyword)
}

func ToProblemFindResponse(ctx context.Context, result service.ProblemSearchResult) (*entity.ProblemFindResponse, error) {
	problems := make([]entity.Problem, len(result.GetResults()))
	for i, p := range result.GetResults() {
		bytes, err := json.Marshal(p.GetProperties(ctx))
		if err != nil {
			return nil, err
		}

		model, err := entity.NewModel(p)
		if err != nil {
			return nil, err
		}

		problems[i] = entity.Problem{
			Model:       model,
			Number:      p.GetNumber(),
			ProblemType: p.GetProblemType(),
			Properties:  bytes,
		}
	}

	return &entity.ProblemFindResponse{
		TotalCount: result.GetTotalCount(),
		Results:    problems,
	}, nil
}

func ToProblemFindAllResponse(ctx context.Context, result service.ProblemSearchResult) (*entity.ProblemFindAllResponse, error) {
	problems := make([]entity.SimpleProblem, len(result.GetResults()))
	for i, p := range result.GetResults() {
		bytes, err := json.Marshal(p.GetProperties(ctx))
		if err != nil {
			return nil, err
		}

		model, err := entity.NewModel(p)
		if err != nil {
			return nil, err
		}

		problems[i] = entity.SimpleProblem{
			ID:          model.ID,
			Number:      p.GetNumber(),
			ProblemType: p.GetProblemType(),
			Properties:  bytes,
		}
	}

	return &entity.ProblemFindAllResponse{
		TotalCount: result.GetTotalCount(),
		Results:    problems,
	}, nil
}

func ToProblemResponse(ctx context.Context, problem domain.ProblemModel) (*entity.Problem, error) {
	bytes, err := json.Marshal(problem.GetProperties(ctx))
	if err != nil {
		return nil, err
	}

	model, err := entity.NewModel(problem)
	if err != nil {
		return nil, err
	}

	return &entity.Problem{
		Model:       model,
		Number:      problem.GetNumber(),
		ProblemType: problem.GetProblemType(),
		Properties:  bytes,
	}, nil
}

func ToProblemIDsCondition(ctx context.Context, param *entity.ProblemIDsParameter, workbookID domain.WorkbookID) (service.ProblemIDsCondition, error) {
	ids := make([]domain.ProblemID, 0)
	for _, id := range param.IDs {
		ids = append(ids, domain.ProblemID(id))
	}
	return service.NewProblemIDsCondition(workbookID, ids)

}

func ToProblemAddParameter(workbookID domain.WorkbookID, param *entity.ProblemAddParameter) (service.ProblemAddParameter, error) {
	var properties map[string]string
	if err := json.Unmarshal(param.Properties, &properties); err != nil {
		return nil, err
	}

	return service.NewProblemAddParameter(workbookID, param.Number, properties)
}

func ToProblemUpdateParameter(param *entity.ProblemUpdateParameter) (service.ProblemUpdateParameter, error) {
	var properties map[string]string
	if err := json.Unmarshal(param.Properties, &properties); err != nil {
		return nil, err
	}

	return service.NewProblemUpdateParameter(param.Number, properties)
}

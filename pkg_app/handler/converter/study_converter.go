package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
)

// func ToStudyResponse(ctx context.Context, recordbook domain.Recordbook) (*entity.RecordbookResponse, error) {
// 	r, err := recordbook.GetResults(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	results := make(map[uint]int)
// 	for k, v := range r {
// 		results[uint(k)] = v
// 	}

// 	return &entity.RecordbookResponse{
// 		ID:      uint(recordbook.GetWorkbookID()),
// 		Results: results,
// 	}, nil
// }

func ToStudyResult(ctx context.Context, param *entity.StudyResultParameter) (*domain.StudyResultParameter, error) {
	return &domain.StudyResultParameter{
		Result: param.Result,
	}, nil
}

func ToProblemWithLevelList(ctx context.Context, problems []domain.ProblemWithLevel) (*entity.ProblemWithLevelList, error) {
	list := make([]entity.ProblemWithLevel, len(problems))
	for i, p := range problems {
		list[i] = entity.ProblemWithLevel{
			ProblemID: uint(p.ProblemID),
			Level:     p.Level,
		}
	}
	return &entity.ProblemWithLevelList{
		Records: list,
	}, nil
}

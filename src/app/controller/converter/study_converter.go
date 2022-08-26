package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-api/src/app/domain"
	libD "github.com/kujilabo/cocotola-api/src/lib/domain"
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

func ToProblemWithLevelList(ctx context.Context, problems []domain.StudyRecordWithProblemID) (*entity.StudyRecords, error) {
	list := make([]*entity.StudyRecord, len(problems))
	for i, p := range problems {
		list[i] = &entity.StudyRecord{
			ProblemID:      uint(p.ProblemID),
			Level:          p.StudyRecord.Level,
			ResultPrev1:    p.StudyRecord.ResultPrev1,
			Memorized:      p.StudyRecord.Memorized,
			LastAnsweredAt: p.StudyRecord.LastAnsweredAt,
		}
	}
	e := &entity.StudyRecords{
		Records: list,
	}
	return e, libD.Validator.Struct(e)
}

func ToIntValue(ctx context.Context, value int) *entity.IntValue {
	return &entity.IntValue{Value: value}
}

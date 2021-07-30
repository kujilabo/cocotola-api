package converter

import (
	"context"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_app/handler/entity"
)

func ToStudyResponse(context context.Context, recordbook domain.Recordbook) (*entity.RecordbookResponse, error) {
	results := make(map[uint]int)
	for k, v := range recordbook.GetResults() {
		results[uint(k)] = v
	}

	return &entity.RecordbookResponse{
		ID:      uint(recordbook.GetID()),
		Results: results,
	}, nil
}

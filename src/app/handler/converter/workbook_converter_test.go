package converter

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kujilabo/cocotola-api/src/app/domain"
	"github.com/kujilabo/cocotola-api/src/app/handler/entity"
	"github.com/kujilabo/cocotola-api/src/app/service"
	serviceM "github.com/kujilabo/cocotola-api/src/app/service/mock"
	userD "github.com/kujilabo/cocotola-api/src/user/domain"
	"github.com/stretchr/testify/require"
)

func TestToWorkbookSearchResponse(t *testing.T) {
	type args struct {
		result service.WorkbookSearchResult
	}

	time1 := time.Now()
	time2 := time1.AddDate(0, 0, 1)
	model11, err := userD.NewModel(1, 2, time1, time2, 3, 4)
	require.NoError(t, err)
	workbookModel11, err := domain.NewWorkbookModel(model11, userD.SpaceID(5), userD.AppUserID(6), userD.NewPrivileges([]userD.RBACAction{domain.PrivilegeRead}), "a", domain.Lang2JA, "problem_type", "question_text", nil)
	require.NoError(t, err)

	workbookSearchResult0 := new(serviceM.WorkbookSearchResult)
	workbookSearchResult0.On("GetResults").Return([]domain.WorkbookModel{})
	workbookSearchResult0.On("GetTotalCount").Return(0)

	workbookSearchResult1 := new(serviceM.WorkbookSearchResult)
	workbookSearchResult1.On("GetResults").Return([]domain.WorkbookModel{workbookModel11})
	workbookSearchResult1.On("GetTotalCount").Return(7)

	tests := []struct {
		name       string
		args       args
		want       *entity.WorkbookSearchResponse
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "has no workbooks",
			args: args{
				result: workbookSearchResult0,
			},
			want: &entity.WorkbookSearchResponse{
				TotalCount: 0,
				Results:    []*entity.Workbook{},
			},
			wantErr: false,
		},
		{
			name: "has one workbook",
			args: args{
				result: workbookSearchResult1,
			},
			want: &entity.WorkbookSearchResponse{
				TotalCount: 7,
				Results: []*entity.Workbook{
					{
						Model: entity.Model{
							ID:        1,
							Version:   2,
							CreatedAt: time1,
							UpdatedAt: time2,
							CreatedBy: 3,
							UpdatedBy: 4,
						},
						Name:         "a",
						Lang2:        "ja",
						ProblemType:  "problem_type",
						QuestionText: "question_text",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToWorkbookSearchResponse(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToWorkbookSearchResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

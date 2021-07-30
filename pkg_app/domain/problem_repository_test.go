package domain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewProblemAddParameter(t *testing.T) {
	m, err := NewProblemAddParameter(WorkbookID(1), 1, "a", map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, m.Number)
}

func TestNewNewProlemParameter(t *testing.T) {
	type args struct {
		workbookID  WorkbookID
		number      int
		problemType string
		properties  map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *ProblemAddParameter
		wantErr bool
	}{
		{
			name: "workbookID is zero",
			args: args{
				workbookID:  WorkbookID(0),
				number:      1,
				problemType: "a",
				properties:  map[string]string{},
			},
			want:    &ProblemAddParameter{WorkbookID(0), 1, "a", map[string]string{}},
			wantErr: true,
		},
		{
			name: "parameters are valid",
			args: args{
				workbookID:  WorkbookID(1),
				number:      1,
				problemType: "a",
				properties:  map[string]string{},
			},
			want:    &ProblemAddParameter{WorkbookID(1), 1, "a", map[string]string{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProblemAddParameter(tt.args.workbookID, tt.args.number, tt.args.problemType, tt.args.properties)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNewProlemParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

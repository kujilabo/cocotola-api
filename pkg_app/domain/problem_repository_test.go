package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProlemAddParameter(t *testing.T) {
	type args struct {
		workbookID WorkbookID
		number     int
		properties map[string]string
	}
	tests := []struct {
		name           string
		args           args
		wantWorkbookID WorkbookID
		wantNumber     int
		wantProperties map[string]string
		wantErr        bool
	}{
		{
			name: "workbookID is zero",
			args: args{
				workbookID: WorkbookID(0),
				number:     1,
				properties: map[string]string{},
			},
			wantWorkbookID: WorkbookID(0),
			wantNumber:     1,
			wantProperties: map[string]string{},
			wantErr:        true,
		},
		{
			name: "parameters are valid",
			args: args{
				workbookID: WorkbookID(1),
				number:     1,
				properties: map[string]string{},
			},
			wantWorkbookID: WorkbookID(1),
			wantNumber:     1,
			wantProperties: map[string]string{},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProblemAddParameter(tt.args.workbookID, tt.args.number, tt.args.properties)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNewProlemParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantWorkbookID, got.GetWorkbookID())
			assert.Equal(t, tt.wantNumber, got.GetNumber())
			// assert.Equal(t, tt.wantProblemTyhpe, got.GetProblemType())
			assert.Equal(t, tt.wantProperties, got.GetProperties())
		})
	}
}

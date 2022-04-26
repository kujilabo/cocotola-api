package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
)

func TestNewProlemAddParameter(t *testing.T) {
	type args struct {
		workbookID domain.WorkbookID
		number     int
		properties map[string]string
	}
	tests := []struct {
		name           string
		args           args
		wantWorkbookID domain.WorkbookID
		wantNumber     int
		wantProperties map[string]string
		wantErr        bool
	}{
		{
			name: "workbookID is zero",
			args: args{
				workbookID: domain.WorkbookID(0),
				number:     1,
				properties: map[string]string{},
			},
			wantWorkbookID: domain.WorkbookID(0),
			wantNumber:     1,
			wantProperties: map[string]string{},
			wantErr:        true,
		},
		{
			name: "parameters are valid",
			args: args{
				workbookID: domain.WorkbookID(1),
				number:     1,
				properties: map[string]string{},
			},
			wantWorkbookID: domain.WorkbookID(1),
			wantNumber:     1,
			wantProperties: map[string]string{},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProblemAddParameter(tt.args.workbookID, tt.args.number, tt.args.properties)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProblemAddParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantWorkbookID, got.GetWorkbookID())
			assert.Equal(t, tt.wantNumber, got.GetNumber())
			assert.Equal(t, tt.wantProperties, got.GetProperties())
		})
	}
}

func TestNewProlemUpdateParameter(t *testing.T) {
	type args struct {
		number     int
		properties map[string]string
	}
	tests := []struct {
		name           string
		args           args
		wantNumber     int
		wantProperties map[string]string
		wantErr        bool
	}{
		{
			name: "parameters are valid",
			args: args{
				number:     1,
				properties: map[string]string{},
			},
			wantNumber:     1,
			wantProperties: map[string]string{},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProblemUpdateParameter(tt.args.number, tt.args.properties)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProblemUpdateParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantNumber, got.GetNumber())
			assert.Equal(t, tt.wantProperties, got.GetProperties())
		})
	}
}

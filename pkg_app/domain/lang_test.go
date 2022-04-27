package domain

import (
	"reflect"
	"testing"
)

func TestNewLang2(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Lang2
		wantErr bool
	}{
		{
			name:    "en",
			args:    "en",
			want:    Lang2EN,
			wantErr: false,
		},
		{
			name:    "ja",
			args:    "ja",
			want:    Lang2JA,
			wantErr: false,
		},
		{
			name:    "empty string",
			args:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLang2(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLang2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLang2() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNewLang3(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		args    string
// 		want    Lang3
// 		wantErr bool
// 	}{
// 		{
// 			name:    "eng",
// 			args:    "eng",
// 			want:    Lang3ENG,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "jpn",
// 			args:    "jpn",
// 			want:    Lang3JPN,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "empty string",
// 			args:    "",
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := NewLang3(tt.args)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewLang3() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewLang3() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

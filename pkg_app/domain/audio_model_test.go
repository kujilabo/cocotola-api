package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAudio(t *testing.T) {
	type args struct {
		id           uint
		lang         Lang2
		text         string
		audioContent string
	}
	tests := []struct {
		name        string
		args        args
		wantID      uint
		wantLang    Lang2
		wantText    string
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid",
			args: args{
				id:           1,
				lang:         Lang2EN,
				text:         "Hello",
				audioContent: "HELLO_CONTENT",
			},
			wantID:      1,
			wantLang:    Lang2EN,
			wantText:    "Hello",
			wantContent: "HELLO_CONTENT",
			wantErr:     false,
		},
		{
			name: "text is empty",
			args: args{
				id:           1,
				lang:         Lang2EN,
				text:         "",
				audioContent: "HELLO_CONTENT",
			},
			wantErr: true,
		},
		{
			name: "content is empty",
			args: args{
				id:           1,
				lang:         Lang2EN,
				text:         "Hello",
				audioContent: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAudioModel(tt.args.id, tt.args.lang, tt.args.text, tt.args.audioContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantID, got.GetID())
				assert.Equal(t, tt.wantLang, got.GetLang())
				assert.Equal(t, tt.wantText, got.GetText())
				assert.Equal(t, tt.wantContent, got.GetContent())
			}
		})
	}
}

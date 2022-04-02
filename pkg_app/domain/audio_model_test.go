package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAudio(t *testing.T) {
	type args struct {
		id           uint
		lang         string
		text         string
		audioContent string
	}
	tests := []struct {
		name        string
		args        args
		wantID      uint
		wantLang    Lang5
		wantText    string
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid",
			args: args{
				id:           1,
				lang:         "en-US",
				text:         "Hello",
				audioContent: "HELLO_CONTENT",
			},
			wantID:      1,
			wantLang:    Lang5ENUS,
			wantText:    "Hello",
			wantContent: "HELLO_CONTENT",
			wantErr:     false,
		},
		{
			name: "text is empty",
			args: args{
				id:           1,
				lang:         "us-US",
				text:         "",
				audioContent: "HELLO_CONTENT",
			},
			wantErr: true,
		},
		{
			name: "content is empty",
			args: args{
				id:           1,
				lang:         "us-US",
				text:         "Hello",
				audioContent: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang, err := NewLang5(tt.args.lang)
			if err != nil {
				panic(err)
			}
			got, err := NewAudioModel(tt.args.id, lang, tt.args.text, tt.args.audioContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAudio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantID, got.GetID())
				assert.Equal(t, tt.wantLang, got.GetLang())
				assert.Equal(t, tt.wantText, got.GetText())
				assert.Equal(t, tt.wantContent, got.GetAudioContent())
			}
		})
	}
}

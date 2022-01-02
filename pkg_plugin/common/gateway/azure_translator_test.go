package gateway

import (
	"context"
	"testing"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/sirupsen/logrus"
)

func Test_azureTranslatorClient_DictionaryLookup(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()
	translatorClient := NewAzureTranslatorClient("")
	type args struct {
		text     string
		fromLang app.Lang2
		toLang   app.Lang2
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "like",
			args: args{
				text:     "like",
				fromLang: app.Lang2EN,
				toLang:   app.Lang2JA,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translatorClient.DictionaryLookup(bg, tt.args.text, tt.args.fromLang, tt.args.toLang)
			if (err != nil) != tt.wantErr {
				t.Errorf("azureTranslatorClient.DictionaryLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logrus.Debugf("got:=%+v\n", got)
			for _, result := range got {
				logrus.Println(result.Target)
			}
		})
	}
}

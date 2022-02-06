package gateway_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
)

func Test_customTranslationRepository_FindByFirstLetter(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	bg := context.Background()
	for _, db := range dbList() {
		type args struct {
			firstLetter string
			lang        app.Lang2
		}
		result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec("delete from custom_translation")
		assert.NoError(t, result.Error)

		book, err := domain.NewTranslation(1, 1, time.Now(), time.Now(), "book", domain.PosNoun, app.Lang2JA, "本", "custom")
		assert.NoError(t, err)

		result = db.Debug().Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(fmt.Sprintf("insert into custom_translation (version,text,pos,lang,translated) values(%d,'%s',%d,'%s','%s')", uint(book.GetVersion()), book.GetText(), int(book.GetPos()), book.GetLang().String(), book.GetTranslated()))
		assert.NoError(t, result.Error)

		tests := []struct {
			name    string
			args    args
			want    []domain.Translation
			wantErr bool
		}{
			{
				name: "found a record",
				args: args{
					firstLetter: "b",
					lang:        app.Lang2JA,
				},
				want: []domain.Translation{
					book,
				},
				wantErr: false,
			},
		}
		r := gateway.NewCustomTranslationRepository(db)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := r.FindByFirstLetter(bg, tt.args.lang, tt.args.firstLetter)
				if (err != nil) != tt.wantErr {
					t.Errorf("customTranslationRepository.FindByFirstLetter() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil {
					assert.Equal(t, len(got), len(tt.want))
					assert.Equal(t, got[0].GetTranslated(), tt.want[0].GetTranslated())
				}

			})
		}
	}
}

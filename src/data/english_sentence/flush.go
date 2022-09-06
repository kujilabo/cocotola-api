package english_sentence

import (
	"context"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	liberrors "github.com/kujilabo/cocotola-api/src/lib/errors"
	"github.com/kujilabo/cocotola-api/src/lib/log"
	pluginEnglishDomain "github.com/kujilabo/cocotola-api/src/plugin/english/domain"
)

func CreateFlushWorkbook(ctx context.Context, studentService appS.Student) error {
	if err := CreateWorkbook(ctx, studentService, "Flush", [][]string{
		{"This is a good book.", "これは良い本です。"},
	}); err != nil {
		return err
	}
	return nil
}

func CreateWorkbook(ctx context.Context, student appS.Student, workbookName string, sentences [][]string) error {
	logger := log.FromContext(ctx)

	workbookProperties := map[string]string{
		"audioEnabled": "false",
	}
	param, err := appS.NewWorkbookAddParameter(pluginEnglishDomain.EnglishSentenceProblemType, workbookName, appD.Lang2JA, "", workbookProperties)
	if err != nil {
		return liberrors.Errorf("NewWorkbookAddParameter. err: %w", err)
	}

	workbookID, err := student.AddWorkbookToPersonalSpace(ctx, param)
	if err != nil {
		return liberrors.Errorf("student.AddWorkbookToPersonalSpace. err: %w", err)
	}

	workbook, err := student.FindWorkbookByID(ctx, workbookID)
	if err != nil {
		return liberrors.Errorf("student.FindWorkbookByID. err: %w", err)
	}

	for i, sentence := range sentences {
		properties := map[string]string{
			"text":       sentence[0],
			"lang2":      "ja",
			"translated": sentence[1],
		}
		param, err := appS.NewProblemAddParameter(workbookID, i+1, properties)
		if err != nil {
			return liberrors.Errorf("NewProblemAddParameter. err: %w", err)
		}

		problemIDs, err := workbook.AddProblem(ctx, student, param)
		if err != nil {
			return liberrors.Errorf("workbook.NewProblemAddParameter. err: %w", err)
		}
		logger.Infof("problemIDs: %v", problemIDs)
	}

	logger.Infof("Example %d", workbookID)
	return nil
}

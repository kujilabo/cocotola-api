package gateway

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"

	"golang.org/x/xerrors"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	common "github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
)

var (
	posPos        = 1
	lenPos        = posPos + 1
	posTranslated = posPos + 1
	lenTranslated = posTranslated + 1
)

type engliushWordProblemAddParameterCSVReader struct {
	workbookID  app.WorkbookID
	problemType string
	reader      *csv.Reader
	num         int
}

func NewEnglishWordProblemAddParameterCSVReader(workbookID app.WorkbookID, problemType string, reader io.Reader) app.ProblemAddParameterIterator {
	return &engliushWordProblemAddParameterCSVReader{
		workbookID:  workbookID,
		problemType: problemType,
		reader:      csv.NewReader(reader),
		num:         1,
	}
}

func (r *engliushWordProblemAddParameterCSVReader) Next() (app.ProblemAddParameter, error) {
	var line []string
	line, err := r.reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, nil
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to reader.Read. err: %w", err)
	}
	if len(line) == 0 {
		return nil, nil
	}

	pos := common.PosOther
	if len(line) >= lenPos {
		posTmp, err := common.ParsePos(line[posPos])
		if err != nil {
			return nil, xerrors.Errorf("failed to ParsePos. err: %w", err)
		}
		pos = posTmp
	}

	translated := ""
	if len(line) >= lenTranslated {
		translated = line[posTranslated]
	}

	properties := map[string]string{
		"lang":       "ja",
		"text":       line[0],
		"translated": translated,
		"pos":        strconv.Itoa(int(pos)),
	}
	param, err := app.NewProblemAddParameter(r.workbookID, r.num, r.problemType, properties)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}

	r.num++
	return param, nil
}

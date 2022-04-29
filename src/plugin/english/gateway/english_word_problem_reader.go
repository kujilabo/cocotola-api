package gateway

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"

	appD "github.com/kujilabo/cocotola-api/src/app/domain"
	appS "github.com/kujilabo/cocotola-api/src/app/service"
	common "github.com/kujilabo/cocotola-api/src/plugin/common/domain"
	"golang.org/x/xerrors"
)

var (
	posPos        = 1
	lenPos        = posPos + 1
	posTranslated = posPos + 1
	lenTranslated = posTranslated + 1
)

type engliushWordProblemAddParameterCSVReader struct {
	workbookID appD.WorkbookID
	// problemType string
	reader *csv.Reader
	num    int
}

func NewEnglishWordProblemAddParameterCSVReader(workbookID appD.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator {
	return &engliushWordProblemAddParameterCSVReader{
		workbookID: workbookID,
		// problemType: problemType,
		reader: csv.NewReader(reader),
		num:    1,
	}
}

func (r *engliushWordProblemAddParameterCSVReader) Next() (appS.ProblemAddParameter, error) {
	var line []string
	line, err := r.reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}
	if err != nil {
		return nil, xerrors.Errorf("failed to reader.Read. err: %w", err)
	}
	if len(line) == 0 {
		return nil, nil
	}
	if len(line[0]) == 0 {
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
		"lang2":      "ja",
		"text":       line[0],
		"translated": translated,
		"pos":        strconv.Itoa(int(pos)),
	}
	param, err := appS.NewProblemAddParameter(r.workbookID, r.num, properties)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewProblemAddParameter. err: %w", err)
	}

	r.num++
	return param, nil
}

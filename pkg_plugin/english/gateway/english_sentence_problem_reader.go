package gateway

import (
	"encoding/csv"
	"errors"
	"io"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	appS "github.com/kujilabo/cocotola-api/pkg_app/service"
)

type englishSentenceProblemAddParameterCSVReader struct {
	workbookID app.WorkbookID
	// problemType string
	reader *csv.Reader
	num    int
}

func NewEnglishSentenceProblemAddParameterCSVReader(workbookID app.WorkbookID, reader io.Reader) appS.ProblemAddParameterIterator {
	return &englishSentenceProblemAddParameterCSVReader{
		workbookID: workbookID,
		// problemType: problemType,
		reader: csv.NewReader(reader),
		num:    1,
	}
}

func (r *englishSentenceProblemAddParameterCSVReader) Next() (appS.ProblemAddParameter, error) {
	var line []string
	line, err := r.reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	properties := map[string]string{
		"lang":       "ja",
		"text":       line[1],
		"translated": line[2],
	}

	param, err := appS.NewProblemAddParameter(r.workbookID, r.num, properties)
	if err != nil {
		return nil, err
	}

	r.num++
	return param, nil
}

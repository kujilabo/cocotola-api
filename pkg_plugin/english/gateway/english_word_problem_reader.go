package gateway

import (
	"encoding/csv"
	"io"
	"strconv"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/domain"
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
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	pos, err := domain.ParsePos(line[0])
	properties := map[string]string{
		"lang":       "ja",
		"text":       line[1],
		"translated": line[2],
		"pos":        strconv.Itoa(int(pos)),
	}
	param, err := app.NewProblemAddParameter(r.workbookID, r.num, r.problemType, properties)
	if err != nil {
		return nil, err
	}

	r.num++
	return param, nil
}

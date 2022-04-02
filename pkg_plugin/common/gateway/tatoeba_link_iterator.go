package gateway

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strconv"

	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
)

type tatoebaLinkAddParameterReader struct {
	reader *csv.Reader
	num    int
}

func NewTatoebaLinkAddParameterReader(reader io.Reader) service.TatoebaLinkAddParameterIterator {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true

	return &tatoebaLinkAddParameterReader{
		reader: csvReader,
		num:    1,
	}
}

func (r *tatoebaLinkAddParameterReader) Next(ctx context.Context) (service.TatoebaLinkAddParameter, error) {
	var line []string
	line, err := r.reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	from, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, err
	}

	to, err := strconv.Atoi(line[1])
	if err != nil {
		return nil, err
	}

	param, err := service.NewTatoebaLinkAddParameter(from, to)
	if err != nil {
		return nil, err
	}

	r.num++
	return param, nil
}

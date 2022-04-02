package gateway

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	app "github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/kujilabo/cocotola-api/pkg_lib/log"
	"github.com/kujilabo/cocotola-api/pkg_plugin/common/service"
	"golang.org/x/xerrors"
)

const (
	bufSize         = 4096
	textLimitLength = 100
)

type tatoebaSentenceAddParameterReader struct {
	// reader *csv.Reader
	reader *bufio.Reader
	num    int
}

// type wrappedReader struct {
// 	reader io.Reader
// }

// func (r *wrappedReader) Read(p []byte) (int, error) {
// 	n, err := r.reader.Read(p)
// 	if err != nil {
// 		return 0, err
// 	}
// 	s1 := string(p)
// 	s2 := strings.ReplaceAll(s1, "\"", "\"\"")
// 	return n + len(s2) - len(s1), nil
// }

func NewTatoebaSentenceAddParameterReader(reader io.Reader) service.TatoebaSentenceAddParameterIterator {
	bufReader := bufio.NewReaderSize(reader, bufSize)
	// wrappedReader:=

	// csvReader := csv.NewReader(reader)
	// csvReader.Comma = '\t'
	// csvReader.LazyQuotes = true

	return &tatoebaSentenceAddParameterReader{
		// reader: csvReader,
		reader: bufReader,
		num:    1,
	}
}

func (r *tatoebaSentenceAddParameterReader) Next(ctx context.Context) (service.TatoebaSentenceAddParameter, error) {
	logger := log.FromContext(ctx)

	b, _, err := r.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		logger.Infof("zero")
	}
	line := strings.Split(string(b), "\t")
	// var line []string
	// line, err := r.reader.Read()
	// if errors.Is(err, io.EOF) {
	// 	return nil, err
	// }

	if err != nil {
		// skip
		logger.Infof("skip rowNumber: %d, line: %v", r.num, line)
		r.num++
		return nil, nil
	}

	sentenceNumber, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, xerrors.Errorf("failed to parse sentenceNumber. rowNumber: %d, value: %s, err: %w", r.num, line[0], err)
	}

	lang, err := app.NewLang3(line[1])
	if err != nil {
		return nil, xerrors.Errorf("failed to NewLang3. rowNumber: %d, value: %s, err: %w", r.num, line[1], err)
	}

	text := line[2]
	author := line[3]

	time1 := line[4]
	time2 := line[5]
	updatedAt := time.Now()

	if len(text) > textLimitLength {
		// skip
		logger.Debugf("skip long text. rowNumber: %d, text: %s", r.num, text)
		r.num++
		return nil, nil
	}

	//\N	2020-02-23 05:07:26
	if r.isValidDatetime(time1) || r.isValidDatetime(time2) {
		var timeS string
		if r.isValidDatetime(time1) {
			timeS = time1
		} else if r.isValidDatetime(time2) {
			timeS = time2
		}

		timeTmp, err := time.Parse("2006-01-02 15:04:05", timeS)
		if err != nil {
			return nil, xerrors.Errorf("failed to Parse. rowNumber: %d, value: %s, err: %w", r.num, timeS, err)
		}
		updatedAt = timeTmp
	}

	param, err := service.NewTatoebaSentenceAddParameter(sentenceNumber, lang, text, author, updatedAt)
	if err != nil {
		return nil, xerrors.Errorf("failed to NewTatoebaSentenceAddParameter. rowNumber: %d, values: %v, err: %w", r.num, line, err)
	}

	r.num++
	return param, nil
}

func (r *tatoebaSentenceAddParameterReader) isValidDatetime(value string) bool {
	return value != "\\N" && value != "0000-00-00 00:00:00"
}

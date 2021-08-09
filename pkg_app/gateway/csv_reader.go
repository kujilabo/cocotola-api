package gateway

import (
	"encoding/csv"
	"io"
)

func ReadCSV(fileReader io.Reader, fn func(i int, line []string) error) error {
	csvReader := csv.NewReader(fileReader)
	var i = 1
	for {
		var line []string
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := fn(i, line); err != nil {
			return err
		}
		i++
	}
	return nil
}

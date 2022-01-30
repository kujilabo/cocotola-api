package main

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/kujilabo/cocotola-api/pkg_plugin/common/gateway"
)

func run() error {
	filePath := "../cocotola-data/datasource/tatoeba/eng_sentences_detailed2.tsv"

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	iterator := gateway.NewTatoebaSentenceAddParameterReader(file)

	for {
		param, err := iterator.Next(context.Background())
		if err != nil {
			return err
		}
		if param == nil {
			continue
		}
	}
}

func main() {
	if err := run(); err != nil {
		if errors.Is(err, io.EOF) {
			return
		}
		panic(err)
	}
}

package domain

import "fmt"

type Lang2 string
type Lang3 string
type Lang5 string

const Lang2Len = 2
const Lang5Len = 5

var (
	Lang2JA Lang2 = "ja"
	Lang2EN Lang2 = "en"

	Lang3JPN Lang3 = "jpn"
	Lang3ENG Lang3 = "eng"

	Lang5ENUS Lang5 = "en-US"
)

func NewLang2(lang string) (Lang2, error) {
	if len(lang) != Lang2Len {
		return Lang2(""), fmt.Errorf("invalid parameter. Lang2: %s", lang)
	}

	return Lang2(lang), nil
}

func NewLang5(lang string) (Lang5, error) {
	if len(lang) != Lang5Len {
		return Lang5(""), fmt.Errorf("invalid parameter. Lang5: %s", lang)
	}

	return Lang5(lang), nil
}

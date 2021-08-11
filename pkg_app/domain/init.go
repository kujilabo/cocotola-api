package domain

var (
	Lang2JA Lang2
	Lang2EN Lang2

	Lang3JPN Lang3
	Lang3ENG Lang3

	Lang5ENUS Lang5
)

func init() {
	var err error

	Lang2EN, err = NewLang2("en")
	if err != nil {
		panic(err)
	}

	Lang2JA, err = NewLang2("ja")
	if err != nil {
		panic(err)
	}

	Lang3JPN, err = NewLang3("jpn")
	if err != nil {
		panic(err)
	}

	Lang3ENG, err = NewLang3("eng")
	if err != nil {
		panic(err)
	}

	Lang5ENUS, err = NewLang5("en-US")
	if err != nil {
		panic(err)
	}
}

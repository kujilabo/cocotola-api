package domain

var (
	Lang2EN      Lang2
	Lang2ES      Lang2
	Lang2JA      Lang2
	Lang2KO      Lang2
	Lang2Unknown Lang2

	// Lang3ENG     Lang3
	// Lang3ESP     Lang3
	// Lang3JPN     Lang3
	// Lang3KOR     Lang3
	// Lang3Unknown Lang3

	// Lang5ENUS    Lang5
	// Lang5ESES    Lang5
	// Lang5JAJP    Lang5
	// Lang5KOKR    Lang5
	// Lang5Unknown Lang5
)

func init() {
	var err error

	Lang2EN, err = NewLang2("en")
	if err != nil {
		panic(err)
	}
	Lang2ES, err = NewLang2("es")
	if err != nil {
		panic(err)
	}
	Lang2JA, err = NewLang2("ja")
	if err != nil {
		panic(err)
	}
	Lang2KO, err = NewLang2("ko")
	if err != nil {
		panic(err)
	}
	Lang2Unknown, err = NewLang2("__")
	if err != nil {
		panic(err)
	}

	// Lang3ENG, err = NewLang3("eng")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang3ESP, err = NewLang3("esp")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang3JPN, err = NewLang3("jpn")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang3KOR, err = NewLang3("kor")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang3Unknown, err = NewLang3("___")
	// if err != nil {
	// 	panic(err)
	// }

	// Lang5ENUS, err = NewLang5("en-US")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang5ESES, err = NewLang5("es-ES")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang5JAJP, err = NewLang5("ja-JP")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang5KOKR, err = NewLang5("ko-KR")
	// if err != nil {
	// 	panic(err)
	// }
	// Lang5Unknown, err = NewLang5("_____")
	// if err != nil {
	// 	panic(err)
	// }
}

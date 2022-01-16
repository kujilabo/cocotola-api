package entity

type TranslationFindParameter struct {
	Letter string `json:"letter"`
}

type Translation struct {
	Lang       string `json:"lang"`
	Text       string `json:"text"`
	Pos        int    `json:"pos"`
	Translated string `json:"translated"`
	Provider   string `json:"provider"`
}

type TranslationFindResponse struct {
	Results []Translation `json:"results"`
}

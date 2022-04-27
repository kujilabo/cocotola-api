package entity

type TranslationFindParameter struct {
	Letter string `json:"letter"`
	Lang2  string `json:"lang2"`
}

type Translation struct {
	Lang2      string `json:"lang2"`
	Text       string `json:"text"`
	Pos        int    `json:"pos"`
	Translated string `json:"translated"`
	Provider   string `json:"provider"`
}

type TranslationFindResponse struct {
	Results []Translation `json:"results"`
}

type TranslationAddParameter struct {
	Lang2      string `json:"lang2" binding:"required"`
	Text       string `json:"text" binding:"required"`
	Pos        int    `json:"pos" binding:"required"`
	Translated string `json:"translated" binding:"required"`
}

type TranslationUpdateParameter struct {
	Translated string `json:"translated" binding:"required"`
}

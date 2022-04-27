package entity

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang2   string `json:"lang2"`
	Text    string `json:"text"`
	Content string `json:"content"`
}

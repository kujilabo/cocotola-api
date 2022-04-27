package entity

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang    string `json:"lang"`
	Text    string `json:"text"`
	Content string `json:"content"`
}

package entity

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang2   string `json:"lang2" validate:"len=2"`
	Text    string `json:"text" validate:"required"`
	Content string `json:"content" validate:"required"`
}

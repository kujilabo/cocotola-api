package entity

type Audio struct {
	ID           uint   `validate:"required,gte=1" json:"id"`
	Lang         string `validate:"required,len=5" json:"lang"`
	Text         string `validate:"required" json:"text"`
	AudioContent string `validate:"required" json:"audioContent"`
}

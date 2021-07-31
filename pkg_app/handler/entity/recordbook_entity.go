package entity

type RecordbookResponse struct {
	ID      uint         `json:"id"`
	Results map[uint]int `json:"results"`
}

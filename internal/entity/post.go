package entity

type Post struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Creator string `json:"creator"` // UUID
}

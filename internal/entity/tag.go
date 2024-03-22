package entity

type Tag struct {
	ID      int    `json:"id"`
	Tag     string `json:"tag"`
	PostID  string `json:"post_id"` // UUID
}

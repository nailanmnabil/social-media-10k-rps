package entity

type Comment struct {
	ID      int    `json:"id"`
	PostID  string `json:"post_id"` // UUID
	Comment string `json:"comment"`
	UserID  string `json:"user_id"` // UUID
}

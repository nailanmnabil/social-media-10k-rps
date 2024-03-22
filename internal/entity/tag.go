package entity

type Tag struct {
	ID        int    `json:"id"`
	ProductID string `json:"product_id"`
	Tag       string `json:"tag"`
}

package entity

// TODO add created_at
type Product struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	ImageURL      string `json:"image_url"`
	Stock         int    `json:"stock"`
	Condition     string `json:"condition"`
	IsPurchasable bool   `json:"is_purchasable"`

	UserID string `json:"user_id"`
}

package entity

type Payment struct {
	ID                   string `json:"id"`
	BankAccountID        string `json:"bank_account_id"`
	PaymentProofImageURL string `json:"payment_proof_image_url"`
	Quantity             int    `json:"quantity"`

	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
}

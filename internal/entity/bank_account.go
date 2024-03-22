package entity

type BankAccount struct {
	ID                string `json:"id"`
	BankName          string `json:"bank_name"`
	BankAccountName   string `json:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number"`
	
	UserID            string `json:"user_id"`
}

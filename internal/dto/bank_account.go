package dto

import "github.com/vandenbill/marketplace-10k-rps/internal/entity"

type (
	ReqCreateBankAccount struct {
		BankName          string `json:"bankName" validate:"required,min=5,max=15"`
		BankAccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
		BankAccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=15"`
	}
	ResGetBankAccount struct {
		BankAccountID     string `json:"bankAccountId"`
		BankName          string `json:"bankName"`
		BankAccountName   string `json:"bankAccountName"`
		BankAccountNumber string `json:"bankAccountNumber"`
	}
)

func (d *ResGetBankAccount) ToDto(bank entity.BankAccount) {
	d.BankAccountID = bank.ID
	d.BankAccountName = bank.BankAccountName
	d.BankAccountNumber = bank.BankAccountNumber
	d.BankName = bank.BankName
}

func (d *ReqCreateBankAccount) ToBankAccountEntity(userID string) entity.BankAccount {
	return entity.BankAccount{BankName: d.BankName, BankAccountName: d.BankAccountName,
		BankAccountNumber: d.BankAccountNumber, UserID: userID}
}

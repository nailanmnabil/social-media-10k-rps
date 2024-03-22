package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/marketplace-10k-rps/internal/cfg"
	"github.com/vandenbill/marketplace-10k-rps/internal/dto"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
	"github.com/vandenbill/marketplace-10k-rps/internal/repo"
)

type BankAccountService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newBankAccountService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *BankAccountService {
	return &BankAccountService{repo, validator, cfg}
}

func (u *BankAccountService) Create(ctx context.Context, body dto.ReqCreateBankAccount, userID string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	bankAccount := body.ToBankAccountEntity(userID)
	err = u.repo.BankAccount.Insert(ctx, bankAccount)
	if err != nil {
		return err
	}

	return nil
}

func (u *BankAccountService) Delete(ctx context.Context, bankID string, sub string) error {
	if bankID == "" {
		return ierr.ErrBadRequest
	}
	bankAccount, err := u.repo.BankAccount.GetByID(ctx, bankID)
	if err != nil {
		return err
	}
	if bankAccount.UserID != sub {
		return ierr.ErrForbidden
	}
	err = u.repo.BankAccount.Delete(ctx, bankID)
	if err != nil {
		return err
	}

	return nil
}

func (u *BankAccountService) Get(ctx context.Context, userID string) ([]dto.ResGetBankAccount, error) {
	results := make([]dto.ResGetBankAccount, 0, 10)
	bankAccounts, err := u.repo.BankAccount.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, v := range bankAccounts {
		res := dto.ResGetBankAccount{}
		res.ToDto(v)
		results = append(results, res)
	}

	return results, nil
}

func (u *BankAccountService) Update(ctx context.Context, body dto.ReqCreateBankAccount, bankID, userID string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	bank, err := u.repo.BankAccount.GetByID(ctx, bankID)
	if err != nil {
		return err
	}
	if bank.UserID != userID {
		return ierr.ErrForbidden
	}
	err = u.repo.BankAccount.Delete(ctx, bankID)
	if err != nil {
		return err
	}
	err = u.Create(ctx, body, userID)

	return err
}

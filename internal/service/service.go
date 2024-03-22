package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/vandenbill/marketplace-10k-rps/internal/cfg"
	"github.com/vandenbill/marketplace-10k-rps/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User        *UserService
	Product     *ProductService
	BankAccount *BankAccountService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Product = newProductService(repo, validator, cfg)
	service.BankAccount = newBankAccountService(repo, validator, cfg)

	return &service
}

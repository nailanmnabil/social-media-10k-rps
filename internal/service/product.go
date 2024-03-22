package service

import (
	"context"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/marketplace-10k-rps/internal/cfg"
	"github.com/vandenbill/marketplace-10k-rps/internal/dto"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
	"github.com/vandenbill/marketplace-10k-rps/internal/repo"
	response "github.com/vandenbill/marketplace-10k-rps/pkg/resp"
)

type ProductService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newProductService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *ProductService {
	return &ProductService{repo, validator, cfg}
}

func (u *ProductService) Create(ctx context.Context, body dto.ReqCreateProduct, userID string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	product := body.ToProductEntity()
	product.UserID = userID

	productID, err := u.repo.Product.Insert(ctx, product)
	if err != nil {
		return err
	}
	err = u.repo.Tag.BatchInsert(ctx, body.ToTagsEntity(productID))
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductService) Delete(ctx context.Context, productID string, sub string) error {
	if productID == "" {
		return ierr.ErrBadRequest
	}
	userID, err := u.repo.Product.FindUserID(ctx, productID)
	if err != nil {
		return err
	}
	if userID != sub {
		return ierr.ErrForbidden
	}
	err = u.repo.Product.Delete(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductService) Update(ctx context.Context, body dto.ReqUpdateProduct, productID string, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	if productID == "" {
		return ierr.ErrNotFound
	}
	userID, err := u.repo.Product.FindUserID(ctx, productID)
	if err != nil {
		return err
	}
	if userID != sub {
		return ierr.ErrForbidden
	}
	err = u.repo.Product.Update(ctx, body.ToProductEntity(productID))
	if err != nil {
		return err
	}
	err = u.repo.Tag.DeleteByProductID(ctx, productID)
	if err != nil {
		return err
	}

	err = u.repo.Tag.BatchInsert(ctx, body.ToTagsEntity(productID))
	return err
}

func (u *ProductService) ChangeStock(ctx context.Context, body dto.ReqChangeStock, productID string, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	if productID == "" {
		return ierr.ErrBadRequest
	}
	userID, err := u.repo.Product.FindUserID(ctx, productID)
	if err != nil {
		return err
	}
	if userID != sub {
		return ierr.ErrForbidden
	}
	err = u.repo.Product.ChangeStock(ctx, productID, body.Stock)
	return err
}

func (u *ProductService) GetByID(ctx context.Context, productID string) (dto.ResGetProductByID, error) {
	res := dto.ResGetProductByID{}

	if productID == "" {
		return res, ierr.ErrBadRequest
	}
	product, soldTotal, purchaseCount, err := u.repo.Product.FindByIdExtended(ctx, productID)
	if err != nil {
		return res, err
	}
	tags, err := u.repo.Tag.GetAllByProductID(ctx, productID)
	if err != nil {
		return res, err
	}
	bankAccs, err := u.repo.BankAccount.GetAllByUserID(ctx, product.UserID)
	if err != nil {
		return res, err
	}
	name, err := u.repo.User.GetNameByID(ctx, product.UserID)
	if err != nil {
		return res, err
	}

	res.Product.ToDto(product, purchaseCount, tags)
	res.Seller.ToDto(bankAccs, soldTotal, name)
	return res, err
}

func (u *ProductService) GetWithFilter(ctx context.Context, filter dto.SearchProductFilter) ([]dto.ResGetProduct, response.Meta, error) {
	meta := response.Meta{}

	if filter.UserOnly && filter.Sub == "" {
		return nil, meta, ierr.ErrForbidden
	}

	res, err := u.repo.Product.GetWithPage(ctx, filter)
	if err != nil {
		return nil, meta, err
	}

	totalRow, err := u.repo.Product.Count(ctx)
	if err != nil {
		return nil, meta, err
	}

	meta.Limit = filter.Limit
	meta.Offset = filter.Offset
	meta.Total = int(math.Ceil(float64(totalRow) / float64(filter.Limit)))

	return res, meta, err
}

// TODO implement transaction
func (u *ProductService) Buy(ctx context.Context, body dto.ReqBuy, productID string, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}
	if productID == "" {
		return ierr.ErrBadRequest
	}
	_, err = u.repo.BankAccount.GetByID(ctx, body.BankAccountId)
	if err != nil {
		return err
	}
	product, err := u.repo.Product.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.Stock-body.Quantity < 0 {
		return ierr.ErrBadRequest
	}

	err = u.repo.Product.ChangeStock(ctx, productID, product.Stock-body.Quantity)
	if err != nil {
		return err
	}

	err = u.repo.Payment.Buy(ctx, body.ToPaymentEntity(productID, sub))
	return err
}

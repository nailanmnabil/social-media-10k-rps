package dto

import (
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
)

type (
	ReqCreateProduct struct {
		Name          string   `json:"name" validate:"required,min=5,max=60"`
		Price         int      `json:"price" validate:"required,min=0"`
		ImageURL      string   `json:"imageUrl" validate:"required,url"`
		Stock         int      `json:"stock" validate:"required,min=0"`
		Condition     string   `json:"condition" validate:"required,oneof=new second"`
		Tags          []string `json:"tags" validate:"required,dive,required"`
		IsPurchasable bool     `json:"isPurchaseable" validate:"required"`
	}
	ReqUpdateProduct struct {
		Name          string   `json:"name" validate:"required,min=5,max=60"`
		Price         int      `json:"price" validate:"required,min=0"`
		ImageURL      string   `json:"imageUrl" validate:"required,url"`
		Condition     string   `json:"condition" validate:"required,oneof=new second"`
		Tags          []string `json:"tags" validate:"required,dive,required"`
		IsPurchasable bool     `json:"isPurchasable" validate:"required"`
	}
	ReqChangeStock struct {
		Stock int `json:"stock" validate:"required,min=0"`
	}
	ReqBuy struct {
		BankAccountId        string `json:"bankAccountId" validate:"required"`
		PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,url"`
		Quantity             int    `json:"quantity" validate:"required,min=1"`
	}
	ResGetProduct struct {
		ProductId     string   `json:"productId"`
		Name          string   `json:"name"`
		Price         int      `json:"price"`
		ImageUrl      string   `json:"imageUrl"`
		Stock         int      `json:"stock"`
		Condition     string   `json:"condition"`
		Tags          []string `json:"tags"`
		IsPurchasable bool     `json:"isPurchasable"`
		PurchaseCount int      `json:"purchaseCount"`
	}
	ResSeller struct {
		Name             string              `json:"name"`
		ProductSoldTotal int                 `json:"productSoldTotal"`
		BankAccounts     []ResGetBankAccount `json:"bankAccounts"`
	}
	ResGetProductByID struct {
		Product ResGetProduct `json:"product"`
		Seller  ResSeller     `json:"seller"`
	}
	SearchProductFilter struct {
		Search         string   `json:"search"`
		Tags           []string `json:"tags"`
		UserOnly       bool     `json:"userOnly"`
		Sub            string   `json:"sub"`
		Condition      string   `json:"condition"`
		ShowEmptyStock bool     `json:"showEmptyStock"`
		MaxPrice       int      `json:"maxPrice"`
		MinPrice       int      `json:"minPrice"`
		SortBy         string   `json:"sortBy"`
		OrderBy        string   `json:"orderBy"`
		Limit          int      `json:"limit"`
		Offset         int      `json:"offset"`
	}
)

func (d *SearchProductFilter) SetDefault() {
	if d.Limit == 0 {
		d.Limit = 10
	}
	if d.Offset == 0 {
		d.Offset = 0
	}
}

func (d *ResGetProduct) ToDto(product entity.Product, purchaseCount int, tags []string) {
	d.ProductId = product.ID
	d.Name = product.Name
	d.Price = product.Price
	d.ImageUrl = product.ImageURL
	d.Stock = product.Stock
	d.Condition = product.Condition
	d.Tags = tags
	d.IsPurchasable = product.IsPurchasable
	d.PurchaseCount = purchaseCount
}

func (d *ResSeller) ToDto(bankAccounts []entity.BankAccount, soldTotal int, name string) {
	d.Name = name
	d.ProductSoldTotal = soldTotal

	resBankAccs := make([]ResGetBankAccount, 0, 10)
	for _, v := range bankAccounts {
		resBankAcc := ResGetBankAccount{BankAccountID: v.ID, BankName: v.BankName,
			BankAccountName: v.BankAccountName, BankAccountNumber: v.BankAccountNumber}
		resBankAccs = append(resBankAccs, resBankAcc)
	}

	d.BankAccounts = resBankAccs
}

func (d *ReqBuy) ToPaymentEntity(productID, userID string) entity.Payment {
	return entity.Payment{BankAccountID: d.BankAccountId, PaymentProofImageURL: d.PaymentProofImageUrl,
		Quantity: d.Quantity, ProductID: productID, UserID: userID}
}

func (d *ReqCreateProduct) ToProductEntity() entity.Product {
	return entity.Product{Name: d.Name, Price: d.Price, ImageURL: d.ImageURL, Stock: d.Stock,
		Condition: d.Condition, IsPurchasable: d.IsPurchasable}
}

func (d *ReqCreateProduct) ToTagsEntity(productID string) []entity.Tag {
	tags := make([]entity.Tag, 0, len(d.Tags))
	for _, v := range d.Tags {
		tag := entity.Tag{ProductID: productID, Tag: v}
		tags = append(tags, tag)
	}

	return tags
}

func (d *ReqUpdateProduct) ToProductEntity(productID string) entity.Product {
	return entity.Product{ID: productID, Name: d.Name, Price: d.Price, ImageURL: d.ImageURL,
		Condition: d.Condition, IsPurchasable: d.IsPurchasable}
}

func (d *ReqUpdateProduct) ToTagsEntity(productID string) []entity.Tag {
	tags := make([]entity.Tag, 0, len(d.Tags))
	for _, v := range d.Tags {
		tag := entity.Tag{ProductID: productID, Tag: v}
		tags = append(tags, tag)
	}

	return tags
}

package dto

import (
	"github.com/vandenbill/social-media-10k-rps/internal/entity"
	"github.com/vandenbill/social-media-10k-rps/pkg/auth"
	"github.com/vandenbill/social-media-10k-rps/pkg/validator"
)

type (
	ReqRegister struct {
		CredentialType  validator.CredentialType `json:"credentialType" validate:"required,oneof=phone email"`
		CredentialValue string                   `json:"credentialValue" validate:"required"`
		Name            string                   `json:"name" validate:"required,min=5,max=50"`
		Password        string                   `json:"password" validate:"required,min=5,max=15"`
	}
	ResRegister struct {
		Phone       string `json:"phone,omitempty"`
		Email       string `json:"email,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqLogin struct {
		CredentialType  validator.CredentialType `json:"credentialType" validate:"required,oneof=phone email"`
		CredentialValue string                   `json:"credentialValue" validate:"required"`
		Password        string                   `json:"password" validate:"required,min=5,max=15"`
	}
	ResLogin struct {
		Phone       string `json:"phone,omitempty"`
		Email       string `json:"email,omitempty"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqLinkEmail struct {
		Email string `json:"email" validate:"required,email"`
	}
	ReqLinkPhone struct {
		Phone string `json:"phone" validate:"required,e164"`
	}
	ReqUpdateAccount struct {
		ImageURL string `json:"imageUrl" validate:"required,url"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
	}
)

func (d *ReqRegister) ToEntity(cryptCost int) (bool, entity.User) {
	email := ""
	phone := ""
	isUseEmail := false

	if d.CredentialType == validator.EmailType {
		email = d.CredentialValue
		isUseEmail = true
	}
	if d.CredentialType == validator.PhoneType {
		phone = d.CredentialValue
	}

	return isUseEmail, entity.User{Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost), Email: email, PhoneNumber: phone}
}

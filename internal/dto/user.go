package dto

import (
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
	"github.com/vandenbill/marketplace-10k-rps/pkg/auth"
)

type (
	ReqRegister struct {
		Username string `json:"username" validate:"required,min=5,max=15"`
		Name     string `json:"name" validate:"required,min=5,max=50"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}
	ResRegister struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	ReqLogin struct {
		Username string `json:"username" validate:"required,min=5,max=15"`
		Password string `json:"password" validate:"required,min=5,max=15"`
	}
	ResLogin struct {
		Username    string `json:"username"`
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
)

func (d *ReqRegister) ToEntity(cryptCost int) entity.User {
	return entity.User{Username: d.Username, Name: d.Name, Password: auth.HashPassword(d.Password, cryptCost)}
}

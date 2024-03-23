package service

import (
	"context"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	"github.com/vandenbill/social-media-10k-rps/internal/repo"
	"github.com/vandenbill/social-media-10k-rps/pkg/auth"
	validatorPkg "github.com/vandenbill/social-media-10k-rps/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) Register(ctx context.Context, body dto.ReqRegister) (dto.ResRegister, error) {
	res := dto.ResRegister{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	if body.CredentialType == validatorPkg.EmailType {
		_, err := mail.ParseAddress(body.CredentialValue)
		if err != nil {
			return res, ierr.ErrBadRequest
		}
	} else if body.CredentialType == validatorPkg.PhoneType {
		v := struct {
			Phone string `validate:"required,e164"`
		}{Phone: body.CredentialValue}

		err := u.validator.Struct(v)
		if err != nil {
			return res, ierr.ErrBadRequest
		}
	}

	isUseEmail, user := body.ToEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user, isUseEmail)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		return res, err
	}

	if isUseEmail {
		res.Email = body.CredentialValue
	} else {
		res.Phone = body.CredentialValue
	}
	res.Name = body.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) Login(ctx context.Context, body dto.ReqLogin) (dto.ResLogin, error) {
	res := dto.ResLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		return res, ierr.ErrBadRequest
	}

	isUseEmail := true
	if body.CredentialType == validatorPkg.EmailType {
		_, err := mail.ParseAddress(body.CredentialValue)
		if err != nil {
			return res, ierr.ErrBadRequest
		}
	} else if body.CredentialType == validatorPkg.PhoneType {
		isUseEmail = false
		v := struct {
			Phone string `validate:"required,e164"`
		}{Phone: body.CredentialValue}

		err := u.validator.Struct(v)
		if err != nil {
			return res, ierr.ErrBadRequest
		}
	}

	user, err := u.repo.User.GetByEmailOrPhone(ctx, body.CredentialValue, isUseEmail)
	if err != nil {
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.Email = user.Email
	res.Phone = user.PhoneNumber
	res.Name = user.Name
	res.AccessToken = token

	return res, nil
}

func (u *UserService) LinkEmail(ctx context.Context, body dto.ReqLinkEmail, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	user, err := u.repo.User.GetByID(ctx, sub)
	if err != nil {
		return err
	}
	if user.Email != "" {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LinkEmail(ctx, body.Email, sub)
	return err
}

func (u *UserService) LinkPhone(ctx context.Context, body dto.ReqLinkPhone, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	user, err := u.repo.User.GetByID(ctx, sub)
	if err != nil {
		return err
	}
	if user.PhoneNumber != "" {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LinkPhone(ctx, body.Phone, sub)
	return err
}

func (u *UserService) UpdateAccount(ctx context.Context, body dto.ReqUpdateAccount, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	if body.ImageURL == "http://incomplete"{
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LookUp(ctx, sub)
	if err != nil {
		return err
	}

	err = u.repo.User.UpdateAccount(ctx, sub, body.Name, body.ImageURL)
	return err
}

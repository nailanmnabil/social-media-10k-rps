package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	"github.com/vandenbill/social-media-10k-rps/internal/repo"
)

type FriendService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newFriendService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *FriendService {
	return &FriendService{repo, validator, cfg}
}

func (u *FriendService) AddFriend(ctx context.Context, body dto.ReqAddFriend, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	if body.UserID == sub {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LookUp(ctx, body.UserID)
	if err != nil {
		return err
	}

	err = u.repo.Friend.AddFriend(ctx, sub, body.UserID)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

func (u *FriendService) DeleteFriend(ctx context.Context, body dto.ReqDeleteFriend, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	if body.UserID == sub {
		return ierr.ErrBadRequest
	}

	err = u.repo.User.LookUp(ctx, body.UserID)
	if err != nil {
		return err
	}

	err = u.repo.Friend.FindFriend(ctx, sub, body.UserID)
	if err != nil {
		if err == ierr.ErrNotFound {
			return ierr.ErrBadRequest
		}
		return err
	}

	err = u.repo.Friend.DeleteFriend(ctx, sub, body.UserID)
	if err != nil {
		return err
	}

	return nil
}

// func (u *FriendService) GetFriends(ctx context.Context, param dto.ParamGetFriends, sub string) error {
// 	err := u.validator.Struct(param)
// 	if err != nil {
// 		return ierr.ErrBadRequest
// 	}

// 	res, err := u.repo.Friend.GetFriends(ctx, param, sub)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

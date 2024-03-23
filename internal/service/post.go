package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	"github.com/vandenbill/social-media-10k-rps/internal/repo"
)

type PostService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newPostService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *PostService {
	return &PostService{repo, validator, cfg}
}

func (u *PostService) AddPost(ctx context.Context, body dto.ReqAddPost, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	isHaveFriend, err := u.repo.Post.IsHaveFriend(ctx, sub)
	if err != nil {
		return err
	}
	if !isHaveFriend {
		return ierr.ErrBadRequest
	}

	postID, err := u.repo.Post.AddPost(ctx, sub, body.PostInHTML)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	err = u.repo.Tag.BatchInsert(ctx, body.Tags, postID)
	return err
}

func (u *PostService) AddComment(ctx context.Context, body dto.ReqAddComment, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		return ierr.ErrBadRequest
	}

	creatorID, err := u.repo.Post.FindPostCreator(ctx, body.PostID)
	if err != nil {
		return err
	}

	err = u.repo.Friend.FindFriend(ctx, sub, creatorID)
	if err != nil {
		if err == ierr.ErrNotFound {
			return ierr.ErrBadRequest
		}
		return err
	}

	err = u.repo.Post.AddComment(ctx, sub, body.PostID, body.Comment)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

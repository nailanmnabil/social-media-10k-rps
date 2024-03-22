package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User   *UserService
	Friend *FriendService
	Post   *PostService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Friend = newFriendService(repo, validator, cfg)
	service.Post = newPostService(repo, validator, cfg)

	return &service
}

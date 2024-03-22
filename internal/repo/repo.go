package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User   *userRepo
	Tag    *tagRepo
	Friend *friendRepo
	Post   *postRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Tag = newTagRepo(conn)
	repo.Friend = newFriendRepo(conn)
	repo.Post = newPostRepo(conn)

	return &repo
}

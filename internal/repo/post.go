package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
)

type postRepo struct {
	conn *pgxpool.Pool
}

func newPostRepo(conn *pgxpool.Pool) *postRepo {
	return &postRepo{conn}
}

func (u *postRepo) IsHaveFriend(ctx context.Context, sub string) (bool, error) {
	q := `SELECT COUNT(*) AS c FROM friends f WHERE f.a = $1`

	c := 0
	err := u.conn.QueryRow(ctx, q,
		sub).Scan(&c)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, ierr.ErrNotFound
		}
		return false, err
	}

	return c > 0, nil
}

func (u *postRepo) AddPost(ctx context.Context, sub, content string) (string, error) {
	q := `INSERT INTO posts (id, creator, content)
	VALUES (gen_random_uuid(), $1, $2) RETURNING id`

	postID := ""
	err := u.conn.QueryRow(ctx, q,
		sub, content).Scan(&postID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return postID, ierr.ErrDuplicate
			}
		}
		return postID, err
	}

	return postID, nil
}

func (u *postRepo) AddComment(ctx context.Context, sub, postID, comment string) error {
	q := `INSERT INTO comments (user_id, post_id, comment)
	VALUES ($1, $2, $3)`

	_, err := u.conn.Exec(ctx, q,
		sub, postID, comment)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

func (u *postRepo) FindPostCreator(ctx context.Context, id string) (string, error) {
	q := `SELECT creator FROM posts WHERE id = $1`

	creator := ""
	err := u.conn.QueryRow(ctx, q,
		id).Scan(&creator)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		return "", err
	}

	return creator, nil
}

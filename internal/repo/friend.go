package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
)

type friendRepo struct {
	conn *pgxpool.Pool
}

func newFriendRepo(conn *pgxpool.Pool) *friendRepo {
	return &friendRepo{conn}
}

func (u *friendRepo) AddFriend(ctx context.Context, sub, friendSub string) error {
	q := `INSERT INTO friends (a, b)
	VALUES ($1, $2), ($2, $1)`

	_, err := u.conn.Exec(ctx, q,
		sub, friendSub)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return ierr.ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (u *friendRepo) DeleteFriend(ctx context.Context, sub, friendSub string) error {
	q := `DELETE FROM friends WHERE (a = $1 and b = $2) or (a = $2 and b = $1)`
	_, err := u.conn.Exec(ctx, q,
		sub, friendSub)

	if err != nil {
		return err
	}

	return nil
}

func (u *friendRepo) FindFriend(ctx context.Context, sub, friendSub string) error {
	q := `SELECT 1 FROM friends WHERE (a = $1 AND b = $2) OR (a = $2 AND b = $1)`

	v := 0
	err := u.conn.QueryRow(ctx, q,
		sub, friendSub).Scan(&v)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		return err
	}

	return nil
}

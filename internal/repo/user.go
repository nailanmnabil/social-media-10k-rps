package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) (string, error) {
	var userID string

	row := u.conn.QueryRow(ctx,
		`INSERT INTO users (id, username, name, password)
		VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id`,
		user.Username, user.Name, user.Password)
	err := row.Scan(&userID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", ierr.ErrDuplicate
			}
		}
		return "", err
	}

	return userID, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	user := entity.User{}

	err := u.conn.QueryRow(ctx,
		`SELECT id, username, name, password FROM users
		WHERE username = $1`,
		username).Scan(&user.ID, &user.Username, &user.Name, &user.Password)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) GetNameByID(ctx context.Context, id string) (string, error) {
	name := ""
	err := u.conn.QueryRow(ctx,
		`SELECT name FROM users
		WHERE id = $1`,
		id).Scan(&name)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return "", ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return "", ierr.ErrNotFound
			}
		}
		return "", err
	}

	return name, nil
}

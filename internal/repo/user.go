package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/social-media-10k-rps/internal/entity"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User, isUseEmail bool) (string, error) {
	credVal := user.Email
	q := `INSERT INTO users (id, name, email, password, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, now()) RETURNING id`
	if !isUseEmail {
		credVal = user.PhoneNumber
		q = `INSERT INTO users (id, name, phone_number, password, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, now()) RETURNING id`
	}

	var userID string
	err := u.conn.QueryRow(ctx, q,
		user.Name, credVal, user.Password).Scan(&userID)

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

func (u *userRepo) GetByEmailOrPhone(ctx context.Context, cred string, isUseEmail bool) (entity.User, error) {
	user := entity.User{}
	q := `SELECT id, name, password FROM users
	WHERE email = $1`
	if !isUseEmail {
		q = `SELECT id, name, password FROM users
		WHERE phone_number = $1`
	}

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

// func (u *userRepo) GetNameByID(ctx context.Context, id string) (string, error) {
// 	name := ""
// 	err := u.conn.QueryRow(ctx,
// 		`SELECT name FROM users
// 		WHERE id = $1`,
// 		id).Scan(&name)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return "", ierr.ErrNotFound
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "22P02" {
// 				return "", ierr.ErrNotFound
// 			}
// 		}
// 		return "", err
// 	}

// 	return name, nil
// }

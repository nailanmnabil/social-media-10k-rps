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

// func (u *friendRepo) GetFriends(ctx context.Context, param dto.ParamGetFriends, sub string) ([]dto.ResGetFriends, error) {
// 	var query strings.Builder

// 	query.WriteString("SELECT id, name, image_url, created_at, (SELECT COUNT(*) FROM friends f2 WHERE f2.a = u.id) as friendCount FROM friends f WHERE 1 = 1")
// 	if param.OnlyFriend {
// 		query.WriteString("SELECT id, name, image_url, created_at, (SELECT COUNT(*) FROM friends f2 WHERE f2.a = f.b) as friendCount FROM users u WHERE 1 = 1 ")
// 	}

// 	if param.OnlyFriend {
// 		query.WriteString(fmt.Sprintf("AND u.id != %s ", sub))
// 	} else {
// 		query.WriteString(fmt.Sprintf("AND f.a = %s ", sub))
// 	}

// 	if param.Search != "" {
// 		query.WriteString(fmt.Sprintf("AND name LIKE '%s' ", fmt.Sprintf("%%%s%%", param.Search)))
// 	}

// 	query.WriteString(fmt.Sprintf("ORDER BY %s %s ", param.SortBy, param.OrderBy))

// 	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

// 	rows, err := u.conn.Query(ctx, query.String())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	results := make([]dto.ResGetFriends, 0, 10)
// 	for rows.Next() {
// 		result := dto.ResGetFriends{}
// 		err := rows.Scan(&result.UserID, &result.Name, &result.ImageURL, &result.CreatedAt, &result.FriendCount)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, result)
// 	}

// 	return res, nil
// }

package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	timepkg "github.com/vandenbill/social-media-10k-rps/pkg/time"
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

func (u *friendRepo) GetFriends(ctx context.Context, param dto.ParamGetFriends, sub string) ([]dto.ResGetFriends, int, error) {
	var query strings.Builder

	if param.OnlyFriend {
		query.WriteString(fmt.Sprintf("SELECT u.id, u.name, u.image_url, u.created_at, (select count(*) from friends f2 where f2.a = f.b) as friendCount from friends f join users u on u.id = f.b WHERE f.a = '%s' ", sub))
	} else {
		query.WriteString("SELECT u.id, u.name, u.image_url, u.created_at, (SELECT COUNT(*) FROM friends f WHERE f.a = u.id) as friendCount FROM users u WHERE 1 = 1 ")
	}

	if param.Search != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Search)))
	}

	if param.SortBy == "createdAt" {
		param.SortBy = "created_at"
	}
	query.WriteString(fmt.Sprintf("ORDER BY %s %s ", param.SortBy, param.OrderBy))

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	rows, err := u.conn.Query(ctx, query.String())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := make([]dto.ResGetFriends, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt time.Time

		result := dto.ResGetFriends{}
		err := rows.Scan(&result.UserID, &result.Name, &imageUrl, &createdAt, &result.FriendCount)
		if err != nil {
			return nil, 0, err
		}

		result.ImageURL = imageUrl.String
		result.CreatedAt = timepkg.TimeToISO8601(createdAt)
		results = append(results, result)
	}

	count, err := u.count(ctx, query.String())
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (u *friendRepo) count(ctx context.Context, q string) (int, error) {
	q = fmt.Sprintf(`SELECT COUNT(*) AS totalRows FROM (%s)`, q)
	count := 0
	err := u.conn.QueryRow(ctx, q).Scan(&count)
	return count, err
}

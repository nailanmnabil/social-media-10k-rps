package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type tagRepo struct {
	conn *pgxpool.Pool
}

func newTagRepo(conn *pgxpool.Pool) *tagRepo {
	return &tagRepo{conn}
}

func (r *tagRepo) BatchInsert(ctx context.Context, tags []string, postID string) error {
	var values []interface{}
	for _, tag := range tags {
		values = append(values, postID, tag)
	}

	query := "INSERT INTO tags (post_id, tag) VALUES "
	var placeholders []string
	for i := 0; i < len(tags); i++ {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
	}
	query += strings.Join(placeholders, ",")

	_, err := r.conn.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

// func (r *tagRepo) DeleteByProductID(ctx context.Context, productID string) error {
// 	_, err := r.conn.Exec(ctx, `
// 	DELETE FROM tags WHERE post_id = $1
// 	`, productID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *tagRepo) GetAllByProductID(ctx context.Context, productID string) ([]string, error) {
// 	rows, err := r.conn.Query(ctx, `
// 	SELECT tag FROM tags WHERE post_id = $1
// 	`, productID)
// 	if err != nil {
// 		if err.Error() == "no rows in result set" {
// 			return nil, ierr.ErrNotFound
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok {
// 			if pgErr.Code == "22P02" {
// 				return nil, ierr.ErrNotFound
// 			}
// 		}
// 		return nil, err
// 	}

// 	tags := make([]string, 0, 10)
// 	for rows.Next() {
// 		tag := ""
// 		if err := rows.Scan(&tag); err != nil {
// 			return nil, err
// 		}
// 		tags = append(tags, tag)
// 	}

// 	return tags, nil
// }

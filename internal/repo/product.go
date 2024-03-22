package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/vandenbill/marketplace-10k-rps/internal/dto"
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
)

type productRepo struct {
	conn *pgxpool.Pool
}

func newProductRepo(conn *pgxpool.Pool) *productRepo {
	return &productRepo{conn}
}

func (r *productRepo) Insert(ctx context.Context, product entity.Product) (string, error) {
	var productID string

	err := r.conn.QueryRow(ctx, `
		INSERT INTO products (id, name, price, image_url, stock, condition, is_purchasable, user_id)
		VALUES (gen_random_uuid() ,$1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		product.Name, product.Price, product.ImageURL, product.Stock, product.Condition,
		product.IsPurchasable, product.UserID).
		Scan(&productID)
	if err != nil {
		return "", err // TODO check postgres error with casting
	}

	return productID, nil
}

func (r *productRepo) FindByID(ctx context.Context, productID string) (entity.Product, error) {
	product := entity.Product{}
	product.ID = productID

	err := r.conn.QueryRow(ctx, `
		SELECT name, price, image_url, stock, condition, is_purchasable, user_id FROM products
		WHERE id = $1`,
		productID).
		Scan(&product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition,
			&product.IsPurchasable, &product.UserID)
	if err != nil {
		return product, err // TODO check postgres error with casting
	}

	return product, nil
}

func (r *productRepo) FindByIdExtended(ctx context.Context, productID string) (entity.Product, int, int, error) {
	productSoldTotal := 0
	purchaseCount := 0
	product := entity.Product{}
	product.ID = productID

	err := r.conn.QueryRow(ctx, `
		SELECT name, price, image_url, stock, condition, is_purchasable, user_id, 
		COALESCE((SELECT SUM(p.quantity) FROM payments p WHERE p.product_id = prod.id), 0) as product_sold_total,
		COALESCE((SELECT COUNT(p.id) FROM payments p WHERE p.product_id = prod.id), 0) as purchase_count
		FROM products prod
		WHERE id = $1`,
		productID).
		Scan(&product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition,
			&product.IsPurchasable, &product.UserID, &productSoldTotal, &purchaseCount)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return product, 0, 0, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return product, 0, 0, ierr.ErrNotFound
			}
		}
		return product, 0, 0, err
	}

	return product, productSoldTotal, purchaseCount, nil
}

func (r *productRepo) FindUserID(ctx context.Context, productID string) (string, error) {
	userID := ""

	err := r.conn.QueryRow(ctx, `
		SELECT user_id FROM products
		WHERE id = $1`,
		productID).
		Scan(&userID)
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

	return userID, nil
}

func (r *productRepo) ChangeStock(ctx context.Context, productID string, stock int) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE products
		SET stock = $1
		WHERE id = $2`,
		stock, productID)
	return err
}

func (r *productRepo) Delete(ctx context.Context, productID string) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM products WHERE id = $1`,
		productID)
	if err != nil {
		return err // TODO check postgres error with casting
	}

	return nil
}

func (r *productRepo) Update(ctx context.Context, product entity.Product) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE products
		SET name = $1, price = $2, image_url = $3, condition = $4, is_purchasable = $5
		WHERE id = $6`,
		product.Name, product.Price, product.ImageURL, product.Condition, product.IsPurchasable,
		product.ID)
	if err != nil {
		return err // TODO check postgres error with casting
	}

	return nil
}

func (r *productRepo) Count(ctx context.Context) (int, error) {
	totalRow := 0
	err := r.conn.QueryRow(ctx, `SELECT COUNT(id) FROM products`).Scan(&totalRow)
	return totalRow, err
}

func (r *productRepo) GetWithPage(ctx context.Context, filter dto.SearchProductFilter) ([]dto.ResGetProduct, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, price, image_url, stock, condition, is_purchasable FROM products WHERE 1=1")

	if filter.Search != "" {
		query.WriteString(fmt.Sprintf(" AND name LIKE '%s'", fmt.Sprintf("%%%s%%", filter.Search)))
	}

	if len(filter.Tags) > 0 {
		if filter.Tags[0] != "" {
			for i := range filter.Tags {
				filter.Tags[i] = "'" + filter.Tags[i] + "'"
			}
			query.WriteString(fmt.Sprintf(" AND id IN (SELECT product_id FROM tags WHERE tag IN (%s))",
				strings.Join(filter.Tags, ",")))
		}
	}

	if filter.UserOnly {
		query.WriteString(fmt.Sprintf(" AND user_id = '%s'", filter.Sub))
	}

	if filter.Condition != "" {
		query.WriteString(fmt.Sprintf(" AND condition = '%s'", filter.Condition))
	}

	if filter.ShowEmptyStock {
		query.WriteString(" AND stock > -1")
	} else {
		query.WriteString(" AND stock > 0")
	}

	if filter.MaxPrice > 0 {
		query.WriteString(fmt.Sprintf(" AND price <= %d", filter.MaxPrice))
	}

	if filter.MinPrice > 0 {
		query.WriteString(fmt.Sprintf(" AND price >= %d", filter.MinPrice))
	}

	if filter.SortBy != "" {
		query.WriteString(fmt.Sprintf(" ORDER BY %s", filter.SortBy))
		if filter.OrderBy != "" {
			query.WriteString(fmt.Sprintf(" %s", filter.OrderBy))
		}
	}

	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, filter.Offset))

	rows, err := r.conn.Query(ctx, query.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	results := make([]dto.ResGetProduct, 0, filter.Limit)

	for rows.Next() {
		res := dto.ResGetProduct{}
		if err := rows.Scan(&res.ProductId, &res.Name, &res.Price, &res.ImageUrl, &res.Stock, &res.Condition,
			&res.IsPurchasable); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		results = append(results, res)
	}

	return results, nil
}

package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
)

type paymentRepo struct {
	conn *pgxpool.Pool
}

func newPaymentRepo(conn *pgxpool.Pool) *paymentRepo {
	return &paymentRepo{conn}
}

func (r *paymentRepo) Buy(ctx context.Context, payment entity.Payment) error {
	_, err := r.conn.Exec(ctx, `INSERT INTO 
	payments (id, user_id, product_id, bank_account_id, payment_proof_image_url, quantity)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`, payment.UserID, payment.ProductID, payment.BankAccountID,
		payment.PaymentProofImageURL, payment.Quantity)

	if err != nil {
		return err // TODO check postgres error with casting
	}

	return nil
}

package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandenbill/marketplace-10k-rps/internal/entity"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
)

type bankAccountRepo struct {
	conn *pgxpool.Pool
}

func newBankAccountRepo(conn *pgxpool.Pool) *bankAccountRepo {
	return &bankAccountRepo{conn}
}

func (u *bankAccountRepo) Insert(ctx context.Context, bankAccount entity.BankAccount) error {
	_, err := u.conn.Exec(ctx,
		`INSERT INTO bank_accounts (id, user_id, bank_name, bank_account_name, bank_account_number)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)`,
		bankAccount.UserID, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber)

	return err
}

func (u *bankAccountRepo) GetAllByUserID(ctx context.Context, userID string) ([]entity.BankAccount, error) {
	bankAccounts := make([]entity.BankAccount, 0, 10)

	rows, err := u.conn.Query(ctx,
		`SELECT id, bank_name, bank_account_name, bank_account_number 
		FROM bank_accounts WHERE user_id = $1`,
		userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return nil, ierr.ErrNotFound
			}
		}
		return nil, err
	}

	for rows.Next() {
		bankAccount := entity.BankAccount{}
		if err := rows.Scan(&bankAccount.ID, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber); err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, bankAccount)
	}

	return bankAccounts, nil
}

func (u *bankAccountRepo) GetByID(ctx context.Context, bankID string) (entity.BankAccount, error) {
	bankAccount := entity.BankAccount{}
	err := u.conn.QueryRow(ctx,
		`SELECT bank_name, bank_account_name, bank_account_number, user_id
		FROM bank_accounts WHERE id = $1`,
		bankID).Scan(&bankAccount.BankName, &bankAccount.BankAccountName,
		&bankAccount.BankAccountNumber, &bankAccount.UserID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return bankAccount, ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return bankAccount, ierr.ErrNotFound
			}
		}
		return bankAccount, err
	}

	return bankAccount, nil
}

func (r *bankAccountRepo) Delete(ctx context.Context, bankID string) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM bank_accounts WHERE id = $1`,
		bankID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return ierr.ErrNotFound
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "22P02" {
				return ierr.ErrNotFound
			}
		}
		return err
	}

	return nil
}

package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User        *userRepo
	Product     *productRepo
	Tag         *tagRepo
	BankAccount *bankAccountRepo
	Payment     *paymentRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Product = newProductRepo(conn)
	repo.Tag = newTagRepo(conn)
	repo.BankAccount = newBankAccountRepo(conn)
	repo.Payment = newPaymentRepo(conn)

	return &repo
}

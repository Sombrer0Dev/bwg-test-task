package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {
	const op = "storage.portgres.New"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: open db connection: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddAccount(currency string) (uuid.UUID, error) {
	const op = "storage.postgres.Add"

	// using uuid here because it's the easiest way to implement unique value
	wallet := uuid.New()

	stmt, err := s.db.Prepare("INSERT INTO account(currency, wallet) VALUES ($1, $2)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(currency, wallet)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return wallet, nil
}

func (s *Storage) Invoice(currency string, wallet uuid.UUID, amount float64) error {
	const op = "storage.postgres.Invoice"

	if ok, err := s.checkCurrency(wallet, currency); !ok {
		if err != nil {
			return fmt.Errorf("%s: checking currency: %w", op, err)
		}
		return fmt.Errorf("%s: currency does not match", op)
	}

	stmt, err := s.db.Prepare(`UPDATE account SET balance = balance + $1 WHERE wallet = $2`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(amount, wallet)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) checkCurrency(wallet uuid.UUID, currency string) (bool, error) {
	const op = "storage.postgres.getCurrency"

	stmt, err := s.db.Prepare(`SELECT (currency=$1) FROM account WHERE wallet = $2`)
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var match bool
	if err := stmt.QueryRow(currency, wallet).Scan(&match); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("%s: unknown wallet: %w", op, err)
		}
		return false, fmt.Errorf("%s: query statement: %w", op, err)
	}

	return match, nil
}

package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/Sombrer0Dev/bwg-test-task/internal/storage"
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

func (s *Storage) AddAccount(currency string) (uuid.UUID, int64, error) {
	const op = "storage.sqlite.SaveURL"

	// using uuid here because it's the easiest way to implement unique value
	wallet := uuid.New()

	stmt, err := s.db.Prepare("INSERT INTO account(currency, wallet) VALUES (?, ?)")
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(currency, wallet)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("%s: failed to get last inserted id: %w", op, err)
	}
	return wallet, id, nil
}

func (s *Storage) GetURL(id int) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE id = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var res string

	err = stmt.QueryRow(id).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return res, nil
}

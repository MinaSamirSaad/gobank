package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
}

type postgresStorage struct {
	db *sql.DB
}

func PostgresStore() (*postgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &postgresStorage{db: db}, nil
}

func (s *postgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *postgresStorage) createAccountTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			first_name TEXT,
			last_name TEXT,
			number serial,
			balance serial,
			created_at TIMESTAMP
		)
	`)
	return err
}

func (s *postgresStorage) CreateAccount(a *Account) error {
	query := `insert into accounts 
	(first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)
	`
	resp, err := s.db.Query(query, a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("the insert statement %+v\n", resp)
	return nil
}

func (s *postgresStorage) DeleteAccount(id int) error {
	return nil
}

func (s *postgresStorage) UpdateAccount(a *Account) error {
	return nil
}

func (s *postgresStorage) GetAccountById(id int) (*Account, error) {
	return nil, nil
}

func (s *postgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		var account Account
		err = rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
}

type postgresStorage struct {
	db *sql.DB
}

// create the db ana initialize it
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

// create Account Table in postgres Database
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

// create Account record in Accounts Table
func (s *postgresStorage) CreateAccount(a *Account) (*Account, error) {
	query := `insert into accounts 
	(first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)
	returning id, first_name, last_name, number, balance, created_at
	`
	resp, err := s.db.Query(query, a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt)
	if err != nil {
		return nil, err
	}
	for resp.Next() {
		return ScanIntoAccount(resp)
	}
	return nil, fmt.Errorf("some thing wrong happen")
}

// Delete Account record in Accounts Table
func (s *postgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query(`Delete FROM accounts where id = $1`, id)
	return err
}

// Update Account record data in Accounts Table
func (s *postgresStorage) UpdateAccount(a *Account) error {
	return nil
}

// Get Account record data in Accounts Table
func (s *postgresStorage) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts where id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

// Get all Accounts record data in Accounts Table
func (s *postgresStorage) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// helper function to help scan account data from rows
func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: account.sql

package db

import (
	"context"
	"database/sql"
)

const createAccount = `-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  ?, ?, ?
)
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = ? LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :execresult
UPDATE accounts SET balance = ? where id = ?
`

type UpdateAccountParams struct {
	Balance int64 `json:"balance"`
	ID      int64 `json:"id"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateAccount, arg.Balance, arg.ID)
}

-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  ?, ?, ?
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts SET balance = ? where id = ?;

-- name: AddAccountBalance :execresult
UPDATE accounts SET balance = balance + sqlc.arg(amount) where id = sqlc.arg(id);
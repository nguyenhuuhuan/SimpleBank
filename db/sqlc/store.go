package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct {
}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")
		res, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		transferId, err := res.LastInsertId()
		if err != nil {
			return err
		}

		result.Transfer, err = q.GetTransfer(ctx, transferId)
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		fromEntryId, err := fromEntry.LastInsertId()
		if err != nil {
			return err
		}
		resFromEntry, err := q.GetEntry(ctx, fromEntryId)
		if err != nil {
			return err
		}
		result.FromEntry = resFromEntry

		fmt.Println(txName, "create transfer 2")
		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		toEntryId, err := toEntry.LastInsertId()
		if err != nil {
			return err
		}
		resToEntry, err := q.GetEntry(ctx, toEntryId)
		if err != nil {
			return err
		}
		result.ToEntry = resToEntry

		fmt.Println(txName, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 1")
		_, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}
		resFromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount = resFromAccount

		fmt.Println(txName, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 2")
		_, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		resToAccount, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount = resToAccount

		return nil
	})

	return result, err
}

func ResultEntry(ctx context.Context, q *Queries, sql sql.Result) (TransferTxResult, error) {
	var result TransferTxResult
	id, err := sql.LastInsertId()
	if err != nil {
		return result, err
	}

	result.Transfer, err = q.GetTransfer(ctx, id)
	if err != nil {
		return result, err
	}
	return result, nil
}

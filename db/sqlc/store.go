package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute database queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// create New Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbError)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id`
	ToAccountId   int64 `json:"to_account_id`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Tranfer     Tranfer `json:"tranfer"`
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

// TranferTx performs a money transfer from 1 account to the other
// It create a transfer record, add account entries and update account balance with a single db transaction

func (store *Store) TranferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create a transfer record
		result.Tranfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// add account entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		// update account's balance
		// minus money in account 1

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     arg.FromAccountId,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}

			//plus money to account 2

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     arg.ToAccountId,
				Amount: arg.Amount,
			})

			if err != nil {
				return err
			}
		} else {
			//plus money to account 2
			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     arg.ToAccountId,
				Amount: arg.Amount,
			})

			if err != nil {
				return err
			}

			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:     arg.FromAccountId,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err
}

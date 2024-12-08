package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all the function to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
	
}

func NewStore(db *sql.DB) *Store{
	return &Store {
		db: db,
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
// TransferTxParams contains the input parameters of the transfer transaction

// TransferTxResult is the result of he transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}
// Money transfer from one account to the other
// Create a transfer record, add account entries, and update accounts' balance within a single database transaction

func (store *Store) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error{
		var err error
		
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}


		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		}) 
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(context.Background(), q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(context.Background(), q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney (
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
)(account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID: accountID1,
	})
	if err != nil {
		return
	}
	
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID: accountID2,
	})
	return

}
package db

import (
	"context"
	"database/sql"
	"fmt"
)

// provides all funcs to exec db queries and transactions
// composition for extending Queries functionality
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

// executes a function within a db transaction
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

// input params of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// performs a money transfer from one account to another
// it creates a transfer record, adds account entries, update accounts' balance within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		// creating a transfer record
		fmt.Println(txName, "create transfer")
		params := CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}

		result.Transfer, err = q.CreateTransfer(ctx, params)
		if err != nil {
			return err
		}

		// adding account entries
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// updating accounts' balances
		// обновляем всегда первым баланс аккунта с меньшим ID
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
        } else {
			result.FromAccount, result.ToAccount, err = addMoney(q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
        }

		return err
	})

	return result, err
}


func addMoney(
    q *Queries,
    accountID1 int64,
    amount1 int64,
    accountID2 int64,
    amount2 int64,
) (account1 Account, account2 Account, err error)  {
	account1, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		Amount: amount1,
		ID: accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		Amount: amount2,
		ID: accountID2,
	})
	if err != nil {
		return
	}

	return account1, account2, err
}
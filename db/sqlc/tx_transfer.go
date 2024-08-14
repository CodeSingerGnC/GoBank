package db

import (
	"context"
	"errors"
)

// TransferTxParam 定义转账操作需要的参数
type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	TransferID 		int64
	FromEntryID		int64
	ToEntryID		int64
	FromAccountID 	int64
	ToAccountID 	int64
}

const (
	// ErrAffectedTooFewRows = "affected too few rows"
	// ErrAffectedTooManyRows = "affected too many rows"
	ErrAffectedRowsNotMatch = "affected rows not match"
)

// - TODO: 账户会出现负值

// TransferTx 用于完成转账操作，将产生一条转账记录，两条账户条目，以及更新两个账户的余额。
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var transfer TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 修改转出账户余额
		transfer.FromAccountID = arg.FromAccountID

		result, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Balance: -arg.Amount,
		})
		if err != nil {
			return err
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New(ErrAffectedRowsNotMatch)
		}

		// 修改转入账户余额
		transfer.ToAccountID = arg.ToAccountID

		result, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Balance: arg.Amount,
		})
		if err != nil {
			return err
		}
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New(ErrAffectedRowsNotMatch)
		}

		// 创建转账记录
		result, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}
		lastInsertId, _ := result.LastInsertId()
		transfer.TransferID = lastInsertId
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New(ErrAffectedRowsNotMatch)
		}

		// 创建转出账户条目
		result, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		lastInsertId, _ = result.LastInsertId()
		transfer.FromEntryID = lastInsertId
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New(ErrAffectedRowsNotMatch)
		}

		// 创建转入账户条目
		result, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		lastInsertId, _ = result.LastInsertId()
		transfer.ToEntryID = lastInsertId
		rowsAffected, _ = result.RowsAffected()
		if rowsAffected != 1 {
			return errors.New(ErrAffectedRowsNotMatch)
		}

		return nil
	})

	return transfer, err
}
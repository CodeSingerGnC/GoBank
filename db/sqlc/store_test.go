package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testdb)

	require.NotEmpty(t, store)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// 并发执行转账事务
	goN := 2
	
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < goN; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < goN; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// 检查 transfer
		transfer, err := store.GetTransfer(context.Background(), result.TransferID)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// 检查 entries
		fromEntry, err := store.GetEntry(context.Background(), result.FromEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry, err := store.GetEntry(context.Background(), result.ToEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		// 检查 accounts
		fromAccount, err := store.GetAccount(context.Background(), result.FromAccountID)
		require.NoError(t, err)
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount, err := store.GetAccount(context.Background(), result.ToAccountID)
		require.NoError(t, err)
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
	}

	// 检查最终 balance 状态
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(goN)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(goN)*amount, updateAccount2.Balance)
}

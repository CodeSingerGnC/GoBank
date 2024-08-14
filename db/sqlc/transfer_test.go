package db

import (
	"context"
	"testing"
	"time"

	// "time"

	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fa Account, ta Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fa.ID,
		ToAccountID: ta.ID,
		Amount: util.RandomMoney(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	
	rowsAffected, _ := result.RowsAffected()
	require.True(t, rowsAffected == 1)
	lastInsertId, _ := result.LastInsertId()
	require.NotEmpty(t, lastInsertId)

	var transfer Transfer
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE id = ? LIMIT 1;`
	err = testdb.QueryRow(query, lastInsertId).Scan(&transfer.ID, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.CreatedAt)
	require.NoError(t, err)
	require.Equal(t, lastInsertId, transfer.ID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fa := createRandomAccount(t)
	ta := createRandomAccount(t)
	createRandomTransfer(t, fa, ta)
}

func TestGetTransfer(t *testing.T) {
	fa := createRandomAccount(t)
	ta := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, fa, ta)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fa := createRandomAccount(t)
	ta := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fa, ta)
	}

	arg := ListTransfersParams{
		FromAccountID: fa.ID,
		ToAccountID: ta.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == fa.ID || transfer.ToAccountID == ta.ID)
	}
}
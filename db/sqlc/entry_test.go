package db

import (
	"context"
	"testing"
	"time"

	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, a Account) Entry {
	arg := CreateEntryParams{
		AccountID: a.ID,
		Amount: util.RandomMoney(),
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	rowsAffected, _ := result.RowsAffected()
	require.True(t, rowsAffected == 1)
	lastInsertId, _ := result.LastInsertId()
	require.NotEmpty(t, lastInsertId)
	
	var entry Entry
	query := `SELECT id, account_id, amount, created_at FROM entries WHERE id = ? LIMIT 1`
	err = testdb.QueryRow(query, lastInsertId).Scan(&entry.ID, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
	require.NoError(t, err)
	require.Equal(t, lastInsertId, entry.ID)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotEmpty(t, entry.CreatedAt)

	return entry 
}

func TestCreateEntry(t *testing.T) {
	a := createRandomAccount(t)
	createRandomEntry(t, a)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
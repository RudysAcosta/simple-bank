package db

import (
	"context"
	"testing"
	"time"

	"github.com/RudysAcosta/simple-bank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// createRandomEntry
func createRandomEntry(t *testing.T, account *Account) (Entry, error) {

	arg := CreateEntryParams{
		AccountID: pgtype.Int8{Int64: account.ID, Valid: true},
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	return entry, err
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry, err := createRandomEntry(t, &account)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, pgtype.Int8{Int64: account.ID, Valid: true})
}

// TestGetEntry
func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1, err := createRandomEntry(t, &account)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry2.CreatedAt.Time, time.Second)
}

// TestListEntriesByAccount
func TestListEntriesByAccount(t *testing.T) {
	account := createRandomAccount(t)
	n := 10
	var entries []Entry

	for i := 0; i < n; i++ {
		entry, err := createRandomEntry(t, &account)
		require.NoError(t, err)
		entries = append(entries, entry)
	}

	arg := ListEntriesByAccountParams{
		AccountID: pgtype.Int8{Int64: account.ID, Valid: true},
		Limit:     int32(n),
		Offset:    0,
	}

	entries2, err := testQueries.ListEntriesByAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries2, n)

	for i := 0; i < n; i++ {
		require.Equal(t, entries[i].ID, entries2[i].ID)
		require.Equal(t, entries[i].AccountID, entries2[i].AccountID)
		require.Equal(t, entries[i].Amount, entries2[i].Amount)
		require.WithinDuration(t, entries[i].CreatedAt.Time, entries2[i].CreatedAt.Time, time.Second)
	}
}

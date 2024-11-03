package db

import (
	"context"
	"testing"
	"time"

	"github.com/RudysAcosta/simple-bank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// createRandomTransfer creates a random transfer between two accounts for testing purposes.
// If no accounts are provided, it creates two random accounts.
// If one account is provided, it creates a random account as the recipient.
// If two accounts are provided, it uses them as the sender and recipient respectively.
// If more than two accounts are provided, it fails the test.
func createRandomTransfer(t *testing.T, accounts ...*Account) Transfer {

	var fromAccount, toAccount Account

	if len(accounts) == 0 {
		fromAccount = createRandomAccount(t)
		toAccount = createRandomAccount(t)
	} else if len(accounts) == 1 {
		fromAccount = *accounts[0]
		toAccount = createRandomAccount(t)
	} else if len(accounts) == 2 {
		fromAccount = *accounts[0]
		toAccount = *accounts[1]
	} else {
		t.Fatal("too many accounts provided")
	}

	arg := CreateTransferParams{
		FromAccountID: pgtype.Int8{Int64: fromAccount.ID, Valid: true},
		ToAccountID:   pgtype.Int8{Int64: toAccount.ID, Valid: true},
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.ID, transfer2.ID)

	require.WithinDuration(t, transfer.CreatedAt.Time, transfer2.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, &fromAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: pgtype.Int8{Int64: fromAccount.ID, Valid: true},
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListTransfersByFromAccount(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, &fromAccount, &toAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: pgtype.Int8{Int64: fromAccount.ID, Valid: true},
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestListTransfersByToAccount(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, &fromAccount, &toAccount)
	}

	arg := ListTransfersByToAccountParams{
		ToAccountID: pgtype.Int8{Int64: toAccount.ID, Valid: true},
		Limit:       5,
		Offset:      5,
	}

	transfers, err := testQueries.ListTransfersByToAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestTotalAmountSentByAccount(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, &fromAccount, &toAccount)
	}

	arg := ListTransfersParams{
		FromAccountID: pgtype.Int8{Int64: fromAccount.ID, Valid: true},
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	var totalReceived int64

	for _, transfer := range transfers {
		totalReceived += int64(transfer.Amount)
	}

	idFromAccount := pgtype.Int8{Int64: fromAccount.ID, Valid: true}

	totalAmount, err := testQueries.TotalAmountSentByAccount(context.Background(), idFromAccount)
	require.NoError(t, err)
	require.Equal(t, totalReceived, totalAmount)
}

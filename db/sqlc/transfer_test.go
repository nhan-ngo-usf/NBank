package db

import (
	"context"
	"testing"
	"time"

	"github.com/nhan-ngo-usf/NBank/db/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(fromAccount Account, toAccount Account, t *testing.T) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: util.RandomInt(0, 150),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.Equal(t, transfer.FromAccountID, fromAccount.ID)
	require.Equal(t, transfer.ToAccountID, toAccount.ID)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	CreateRandomTransfer(fromAccount, toAccount, t)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	transfer1 := CreateRandomTransfer(fromAccount, toAccount, t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.NotZero(t, transfer2.ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.WithinDuration(t, transfer2.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomTransfer(fromAccount, toAccount, t)
	}
	arg := ListTransfersParams {
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
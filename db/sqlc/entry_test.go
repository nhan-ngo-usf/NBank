package db

import (
	"context"
	"testing"
	"time"

	"github.com/nhan-ngo-usf/NBank/db/util"
	"github.com/stretchr/testify/require"
)
func CreateRandomEntry(randomAccount Account, t *testing.T) Entry {
	arg := CreateEntryParams {
		AccountID: randomAccount.ID,
		Amount: util.RandomInt(0,150),
	}
	testEntry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, testEntry)
	require.NotZero(t, testEntry.ID)
	require.NotZero(t, testEntry.CreatedAt)
	require.Equal(t, randomAccount.ID, testEntry.AccountID)
	require.Equal(t, arg.Amount, testEntry.Amount)
	return testEntry
}
func TestCreateEntry(t *testing.T){
	randomAccount := createRandomAccount(t)
	CreateRandomEntry(randomAccount, t)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := CreateRandomEntry(account, t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry2.ID, entry1.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry2.Amount, entry2.Amount)
	require.WithinDuration(t, entry2.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	randomAccount := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		CreateRandomEntry(randomAccount, t)
	}
	arg := ListEntriesParams{
		AccountID: randomAccount.ID,
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)

	}
}
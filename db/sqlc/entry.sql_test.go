package db

import (
	"context"
	"testing"

	"github.com/mnakhaev/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	acc := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	res, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, entry.ID, res.ID)
}

func TestListEntries(t *testing.T) {
	// Create a number of entries
	entry := createRandomEntry(t)

	arg := ListEntriesParams{
		AccountID: entry.AccountID,
		Limit:     5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 1)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

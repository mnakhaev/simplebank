package db

import (
	"context"
	"testing"

	"github.com/mnakhaev/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	res, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, transfer.ID, res.ID)
}

func TestListTransfers(t *testing.T) {
	// Create a number of transfers
	tr := createRandomTransfer(t)

	arg := ListTransfersParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Limit:         5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 1)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

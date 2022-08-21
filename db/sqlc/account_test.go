package db

import (
	"context"
	"fmt"
	"simple_bank/db/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwer(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	// var account *Account
	account, err := testQueries.CreateAccount(context.Background(), arg)
	fmt.Println("huan", account)

	require.NoError(t, err)
	insertedAccountID, _ := account.LastInsertId()
	result, _ := testQueries.GetAccount(context.Background(), insertedAccountID)

	require.NotEmpty(t, result)
	require.Equal(t, arg.Owner, result.Owner)
	require.Equal(t, arg.Balance, result.Balance)
	require.Equal(t, arg.Currency, result.Currency)

	require.NotZero(t, result.ID)
	require.NotNil(t, result.CreatedAt)

	return result
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)

	accountID, err := account2.LastInsertId()
	require.NoError(t, err)

	result, err := testQueries.GetAccount(context.Background(), accountID)
	require.NoError(t, err)

	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, result.ID)
	require.Equal(t, account1.Owner, result.Owner)
	require.Equal(t, account1.Balance, result.Balance)
	require.Equal(t, account1.Currency, result.Currency)
	require.WithinDuration(t, account1.CreatedAt, result.CreatedAt, time.Second)

}

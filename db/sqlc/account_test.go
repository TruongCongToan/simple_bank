package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/TruongCongToan/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	//using testify to ensure err is nil
	require.NoError(t, err)
	require.NotEmpty(t, account)

	//check equal
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	//check account id is auto generated
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

// test get accounts
func TestGetOneAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:     account1.ID,
		Amount: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotNil(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Currency, account2.Currency)
	// require.Equal(t, arg.Amount, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestGetAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

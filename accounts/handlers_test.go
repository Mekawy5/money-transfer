// Package accounts for all accounts logic
package accounts

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mekawy5/money-transfer/internals/appctx"
	"github.com/stretchr/testify/require"
)

func TestShouldLoadAccounts(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	appctx := appctx.Context{
		DBConn: db,
	}

	expectedAccounts := []Account{
		{
			ID:      "3d253e29-8785-464f-8fa0-9e4b57699db9",
			Name:    "Trupe",
			Balance: 87.11,
		},
		{
			ID:      "17f904c1-806f-4252-9103-74e7a5d3e340",
			Name:    "Fivespan",
			Balance: 946.15,
		},
	}

	mock.ExpectPrepare("SELECT id, name, balance FROM accounts")
	mock.ExpectQuery("SELECT id, name, balance FROM accounts").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow("3d253e29-8785-464f-8fa0-9e4b57699db9", "Trupe", 87.11).
				AddRow("17f904c1-806f-4252-9103-74e7a5d3e340", "Fivespan", 946.15))

	accounts, _ := GetAccounts(appctx)

	if !reflect.DeepEqual(accounts, expectedAccounts) {
		t.Errorf("Failed test GetAccounts got %v but expected %v", accounts, expectedAccounts)
	}
}

func TestTransferWillReturnErrorWhenNegativeBalanceOccurs(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	appctx := appctx.Context{
		DBConn: db,
	}

	senderID := "3d253e29-8785-464f-8fa0-9e4b57699db9"
	receiverID := "17f904c1-806f-4252-9103-74e7a5d3e340"
	senderBalance := 87.11
	receiverBalance := 946.15
	amount := 100.0

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, name, balance FROM accounts WHERE id = ?").
		WithArgs(senderID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow(senderID, "Trupe", senderBalance))

	mock.ExpectQuery("SELECT id, name, balance FROM accounts WHERE id = ?").
		WithArgs(receiverID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow(receiverID, "Fivespan", receiverBalance))

	mock.ExpectExec("UPDATE accounts SET balance = ? WHERE id = ?").
		WithArgs(senderBalance-amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE accounts SET balance = ? WHERE id = ?").
		WithArgs(receiverBalance+amount, receiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectRollback()

	r := TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	}
	_, _, err = Transfer(r, appctx)

	require.EqualError(t, err, errors.New("no enough balance").Error())
}

func TestTransferWillSucceedWhenMeetsCriteria(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	appctx := appctx.Context{
		DBConn: db,
	}

	senderID := "3d253e29-8785-464f-8fa0-9e4b57699db9"
	receiverID := "17f904c1-806f-4252-9103-74e7a5d3e340"
	senderBalance := 87.11
	receiverBalance := 946.15
	amount := 10.0

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, name, balance FROM accounts WHERE id = ?").
		WithArgs(senderID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow(senderID, "Trupe", senderBalance))

	mock.ExpectQuery("SELECT id, name, balance FROM accounts WHERE id = ?").
		WithArgs(receiverID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow(receiverID, "Fivespan", receiverBalance))

	mock.ExpectExec("UPDATE accounts SET balance = ? WHERE id = ?").
		WithArgs(senderBalance-amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE accounts SET balance = ? WHERE id = ?").
		WithArgs(receiverBalance+amount, receiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := TransferRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	}

	from, to, err := Transfer(r, appctx)

	require.Equal(t, err, nil)

	require.Equal(t, from.Balance, senderBalance-amount)

	require.Equal(t, to.Balance, receiverBalance+amount)
}

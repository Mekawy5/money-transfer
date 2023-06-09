// Package accounts for all accounts logic
package accounts

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mekawy5/money-transfer/internals/appctx"
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

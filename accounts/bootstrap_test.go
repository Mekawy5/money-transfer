package accounts

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mekawy5/money-transfer/internals/appctx"
	"github.com/stretchr/testify/require"
)

func TestLoadAccountsReturnErrorIfJsonInvalid(t *testing.T) {
	// setup
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	appctx := appctx.Context{
		DBConn: db,
	}

	mock.
		ExpectExec("CREATE TABLE IF NOT EXISTS accounts (id TEXT, name TEXT, balance REAL)").
		WillReturnResult(sqlmock.NewResult(0, 0))

	accountsJson := []byte("invalid json")
	err = LoadAccounts(appctx, accountsJson)

	require.EqualError(t, err, errors.New("invalid character 'i' looking for beginning of value").Error())
}

func TestLoadAccountsSuccessWithValidJson(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	appctx := appctx.Context{
		DBConn: db,
	}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS accounts (id TEXT, name TEXT, balance REAL)").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectPrepare("INSERT INTO accounts (id, name, balance) VALUES (?, ?, ?)").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectPrepare("SELECT id, name, balance FROM accounts")

	mock.ExpectQuery("SELECT id, name, balance FROM accounts").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "balance"}).
				AddRow("3d253e29-8785-464f-8fa0-9e4b57699db9", "Trupe", 87.11).
				AddRow("17f904c1-806f-4252-9103-74e7a5d3e340", "Fivespan", 946.15).
				AddRow("fd796d75-1bcf-4a95-bf1a-f7b296adb79f", "Wikizz", 3708.11))

	accountsJson := []byte(`
		[{"id":"3d253e29-8785-464f-8fa0-9e4b57699db9","name":"Trupe","balance":"87.11"},
		{"id":"17f904c1-806f-4252-9103-74e7a5d3e340","name":"Fivespan","balance":"946.15"},
		{"id":"fd796d75-1bcf-4a95-bf1a-f7b296adb79f","name":"Wikizz","balance":"3708.11"}]
	`)

	_ = LoadAccounts(appctx, accountsJson)

	accounts, _ := GetAccounts(appctx)

	require.Equal(t, 3, len(accounts))

	require.NoError(t, mock.ExpectationsWereMet())
}

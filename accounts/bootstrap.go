// accounts
package accounts

import (
	"encoding/json"
	"fmt"

	"github.com/Mekawy5/money-transfer/internals/appctx"
)

// LoadAccounts loads accounts data from embedded file into created sqlite table.
func LoadAccounts(appctx appctx.Context, accountsJson []byte) error {
	accounts := []Account{}

	_, err := appctx.DBConn.Exec("CREATE TABLE IF NOT EXISTS accounts (id TEXT, name TEXT, balance REAL)")
	if err != nil {
		return err
	}

	err = json.Unmarshal(accountsJson, &accounts)
	if err != nil {
		return err
	}

	stmt, err := appctx.DBConn.Prepare("INSERT INTO accounts (id, name, balance) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	for _, account := range accounts {
		_, err := stmt.Exec(account.ID, account.Name, account.Balance)
		if err != nil {
			return err
		}
	}

	fmt.Printf("loaded %d accounts into the database. \n", len(accounts))
	return nil
}
